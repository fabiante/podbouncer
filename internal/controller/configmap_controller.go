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
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ConfigMapReconciler reconciles a single ConfigMap object.
//
// The watched ConfigMap is used to configure a PodReconcilerConfig object.
type ConfigMapReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Config *PodReconcilerConfig
}

const (
	configMapObjectNamespace = "podbouncer-system" // TODO: Make configurable
	configMapObjectName      = "podbouncer-config" // TODO: Make configurable
)

// TODO: Check if the permissions below are too broad. Maybe there are permissions this controller does not actually need.
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=configmaps/status,verbs=get
// +kubebuilder:rbac:groups=core,resources=configmaps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Ignore objects which should not be reconciled
	if req.NamespacedName.String() != fmt.Sprintf("%s/%s", configMapObjectNamespace, configMapObjectName) {
		return ctrl.Result{}, nil // TODO: Figure out if predicates are the proper way to filter out pods from kube-system namespace
	}

	// Get object
	var config v1.ConfigMap
	if err := r.Get(ctx, req.NamespacedName, &config); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Retrieve config value
	data := config.Data

	if maxPodAgeStr, found := data["maxPodAge"]; !found {
		// Log error but do not requeue - the error must be fixed manually
		err := errors.New("missing maxPodAge property in ConfigMap")
		logger.Error(err, "Configuration will not be updated")
		return ctrl.Result{}, nil
	} else {
		maxPodAge, err := time.ParseDuration(maxPodAgeStr)

		if err != nil {
			// Log error but do not requeue - the error must be fixed manually
			err := fmt.Errorf("invalid maxPodAge property in ConfigMap: %s", maxPodAgeStr)
			logger.Error(err, "Configuration will not be updated")
			return ctrl.Result{}, nil
		}

		oldMaxPodAge := r.Config.MaxPodAge()

		r.Config.SetMaxPodAge(maxPodAge)

		logger.Info("Configuration updated", "newMaxPodAge", maxPodAge, "currentMaxPodAge", oldMaxPodAge)

		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Watches(&v1.ConfigMap{}, &handler.EnqueueRequestForObject{}). // TODO: Filter for the single config map object for this controller
		Named("configmap").
		Complete(r)
}