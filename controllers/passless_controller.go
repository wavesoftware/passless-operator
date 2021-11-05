/*
Copyright 2021.

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

package controllers

import (
	"context"
	"reflect"

	"github.com/wavesoftware/passless-operator/api/v1alpha1"
	"github.com/wavesoftware/passless-operator/pkg/masterpassword"
	"github.com/wavesoftware/passless-operator/pkg/masterpassword/secrets"
	"github.com/wavesoftware/passless-operator/pkg/passless"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PassLessReconciler reconciles a PassLess object.
type PassLessReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=wavesoftware.pl,resources=passlesses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=wavesoftware.pl,resources=passlesses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=wavesoftware.pl,resources=passlesses/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *PassLessReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// Fetch the PassLess instance
	psls := &v1alpha1.PassLess{}
	err := r.Client.Get(ctx, req.NamespacedName, psls)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Define a new Secret object
	secret := psls.CreateSecret(r.generator(ctx))
	hashcode := secrets.Hashcode(secret)

	// Set PassLess instance as the owner and controller
	if err := controllerutil.SetControllerReference(psls, secret, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if this Secret already exists
	found := &corev1.Secret{}
	err = r.Client.Get(ctx, types.NamespacedName{
		Name:      secret.Name,
		Namespace: secret.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		l.Info("Creating a new secret",
			"Hash", hashcode,
		)
		err = r.Client.Create(ctx, secret)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Secret created successfully - don't requeue
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Secret already exists, check if need update
	if needsUpdate(secret, found) {
		l.Info("Updating a secret",
			"Hash", hashcode,
		)
		err = r.Client.Update(ctx, secret)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PassLessReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PassLess{}).
		Complete(r)
}

func (r *PassLessReconciler) generator(ctx context.Context) masterpassword.Generator {
	resolver := &passless.Resolver{
		Client:  r.Client,
		Context: ctx,
	}
	return masterpassword.NewGenerator(resolver)
}

func needsUpdate(secret, found *corev1.Secret) bool {
	return !reflect.DeepEqual(secret.Data, found.Data)
}
