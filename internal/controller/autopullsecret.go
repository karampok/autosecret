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

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// AutosecretReconciler reconciles a Autosecret object
type AutosecretReconciler struct {
	client.Client
}

//+kubebuilder:rbac:groups="core",resources=secrets,verbs=get;list;watch;create
//+kubebuilder:rbac:groups="core",resources=namespaces,verbs=get;list;watch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *AutosecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	ns := v1.Namespace{}

	if err := r.Get(ctx, types.NamespacedName{Name: req.Name}, &ns); err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get ns: %v", err)
	}
	logger.Info("Acting", "name", req.Name)

	s := v1.Secret{}

	if err := r.Get(ctx, types.NamespacedName{Name: "pull-secret", Namespace: "openshift-config"}, &s); err != nil {
		return ctrl.Result{}, fmt.Errorf("could not get secret: %v", err)
	}

	s.ObjectMeta.Namespace = req.Name
	s.ObjectMeta.Name = "assisted-deployment-pull-secret"
	s.ObjectMeta.ResourceVersion = ""

	logger.Info("Copy", "name", s.ObjectMeta)
	err := r.Create(ctx, &s)
	if apierrors.IsAlreadyExists(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("put back get secret: %v", err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutosecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).Named("autosecret").
		For(&v1.Namespace{}).WithEventFilter(createdByZTP()).Complete(r)
}

func createdByZTP() predicate.Predicate {
	// Looks like
	// apiVersion: v1
	// kind: Namespace
	// metadata:
	//   annotations:
	//     ran.openshift.io/ztp-gitops-generated: '{}'
	return predicate.NewPredicateFuncs(func(o client.Object) bool {
		m := o.GetAnnotations()
		_, ok := m["ran.openshift.io/ztp-gitops-generated"]
		return ok
	})
}
