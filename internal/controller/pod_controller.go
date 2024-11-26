/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const maxPodAge = time.Second * 15 // TODO: Add configurable value for max age

// TODO: Check if the rbac configuration below makes sense: Is the group actually fabitee.de !?
// +kubebuilder:rbac:groups=fabitee.de,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fabitee.de,resources=pods/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=fabitee.de,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Ignore pods which should not be reconciled
	if req.Namespace == "kube-system" {
		return ctrl.Result{}, nil // TODO: Figure out if predicates are the proper way to filter out pods from kube-system namespace
	}

	// Retrieve pod
	var pod v1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Ignore pods which not in a phase where they should be deleted
	if !r.shouldDeletePod(&pod) {
		return ctrl.Result{}, nil
	}

	// Ignore pods which have not yet reached the deletion deadline
	// TODO: Refactor into testable function?
	podCreatedAt := pod.GetCreationTimestamp() // TODO: Check for zero value?
	podAge := time.Since(podCreatedAt.Time)
	if podAge < maxPodAge {
		// Pod is not yet read for deletion - calculate the expected time when it can be deleted
		requeueAfter := maxPodAge - podAge + time.Second
		logger.Info("Pod may eventually be deleted", "age", podAge, "phase", pod.Status.Phase, "deletionIn", requeueAfter, "deletionAt", time.Now().Add(requeueAfter))
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	logger.Info("Deleting non-running pod", "phase", pod.Status.Phase)

	if err := r.Delete(ctx, &pod); err != nil {
		// Retry deletion in one minute
		return ctrl.Result{RequeueAfter: time.Minute}, fmt.Errorf("failed to delete pod: %w", err)
	}

	logger.Info("Pod deleted")

	return ctrl.Result{}, nil
}

func (r *PodReconciler) shouldDeletePod(pod *v1.Pod) bool {
	phase := pod.Status.Phase
	return phase == v1.PodPending || phase == v1.PodSucceeded || phase == v1.PodFailed
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(&v1.Pod{}, &handler.EnqueueRequestForObject{}). // TODO: Figure out if predicates are the proper way to filter out pods from kube-system namespace
		Named("pod").
		Complete(r)
}
