/*
Copyright 2023.

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

	metal3api "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

//+kubebuilder:rbac:groups="metal3.io",resources=baremetalhosts,verbs=get;list;watch;create

// Autobmh reconciles ...
type AutobmhsecretReconciler struct {
	client.Client
	fromSecret string
}

func (r *AutobmhsecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	bmh := metal3api.BareMetalHost{}

	if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &bmh); err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get ns: %v", err)
	}

	// apiVersion: v1
	// kind: Secret
	// metadata:
	//   namespace: opensfhit-config
	//   name: bmh-secret
	// data:
	//   password: QURNSU4=
	//   username: QURNSU4=
	// type: Opaque

	s := v1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Name: "bmh-secret", Namespace: "openshift-config"}, &s); err != nil {
		logger.Info("No secret to clone")
		return ctrl.Result{}, nil
	}
	s.ObjectMeta.Namespace = req.Namespace
	s.ObjectMeta.Name = bmh.Spec.BMC.CredentialsName
	s.ObjectMeta.ResourceVersion = ""

	logger.Info("Copy", "name", s.ObjectMeta)
	err := r.Create(ctx, &s)
	if apierrors.IsAlreadyExists(err) {
		logger.Info("Exists do nothing")
		return ctrl.Result{}, nil
	}
	if err != nil {
		logger.Info("Failed to clone")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutobmhsecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).Named("autobmhsecret").
		For(&metal3api.BareMetalHost{}).WithEventFilter(createdByZTP()).Complete(r)
}

func bmhInZTPNamespace() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(o client.Object) bool {
		return true
	})
}
