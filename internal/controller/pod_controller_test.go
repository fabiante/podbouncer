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
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"testing"
)

var _ = Describe("Pod Controller", func() {
	Context("When reconciling a resource", func() {

		It("should successfully reconcile the resource", func() {

			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})

func Test_PodReconcilerShouldDeletePod(t *testing.T) {
	type Test struct {
		Phase    v1.PodPhase
		Expected bool
	}

	tests := []Test{
		{
			Phase:    v1.PodPending,
			Expected: true,
		},
		{
			Phase:    v1.PodSucceeded,
			Expected: true,
		},
		{
			Phase:    v1.PodFailed,
			Expected: true,
		},
		{
			Phase:    v1.PodRunning,
			Expected: false,
		},
		{
			Phase:    v1.PodUnknown,
			Expected: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("returns expected value %d", i), func(t *testing.T) {
			pod := &v1.Pod{Status: v1.PodStatus{Phase: test.Phase}}
			r := &PodReconciler{}
			require.Equal(t, test.Expected, r.shouldDeletePod(pod), "unexpected return")
		})
	}
}
