package controllers

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	faasv1alpha1 "eywa/k8s-operator/api/v1alpha1"
)

// FunctionReconciler reconciles a Function object
type FunctionReconciler struct {
	client.Client
	Log    *log.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=faas.eywa.rekfuki.dev,resources=functions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=faas.eywa.rekfuki.dev,resources=functions/status,verbs=get;update;patch
func (r *FunctionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	r.Log.Infof("function", req.NamespacedName)
	spew.Dump(req)
	_ = context.Background()

	// your logic here

	return ctrl.Result{}, nil
}

func (r *FunctionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&faasv1alpha1.Function{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
