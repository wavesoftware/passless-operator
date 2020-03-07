package passless

import (
	"context"
	"reflect"

	wavesoftwarev1alpha1 "github.com/wavesoftware/passless-operator/pkg/apis/wavesoftware/v1alpha1"
	"github.com/wavesoftware/passless-operator/pkg/masterpassword"
	"github.com/wavesoftware/passless-operator/pkg/masterpassword/secrets"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_passless")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new PassLess Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	cl := mgr.GetClient()
	resolver := newResolver(cl)
	return &ReconcilePassLess{
		client:    cl,
		scheme:    mgr.GetScheme(),
		generator: masterpassword.NewGenerator(resolver),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("passless-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource PassLess
	err = c.Watch(&source.Kind{Type: &wavesoftwarev1alpha1.PassLess{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Secret and requeue the owner PassLess
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wavesoftwarev1alpha1.PassLess{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePassLess implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePassLess{}

// ReconcilePassLess reconciles a PassLess object
type ReconcilePassLess struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	scheme    *runtime.Scheme
	generator masterpassword.Generator
}

// Reconcile reads that state of the cluster for a PassLess object and makes changes based on the state read
// and what is in the PassLess.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePassLess) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PassLess")
	ctx := context.TODO()

	// Fetch the PassLess instance
	passless := &wavesoftwarev1alpha1.PassLess{}
	err := r.client.Get(ctx, request.NamespacedName, passless)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Secret object
	secret := passless.CreateSecret(r.generator)
	hashcode := secrets.Hashcode(secret)

	// Set PassLess instance as the owner and controller
	if err := controllerutil.SetControllerReference(passless, secret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Secret already exists
	found := &corev1.Secret{}
	err = r.client.Get(ctx, types.NamespacedName{
		Name:      secret.Name,
		Namespace: secret.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new secret",
			"Hash", hashcode,
		)
		err = r.client.Create(ctx, secret)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Secret created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Secret already exists, check if need update
	if needsUpdate(secret, found) {
		reqLogger.Info("Updating a secret",
			"Hash", hashcode,
		)
		err := r.client.Update(ctx, secret)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

func needsUpdate(secret, found *corev1.Secret) bool {
	return !reflect.DeepEqual(secret.Data, found.Data)
}
