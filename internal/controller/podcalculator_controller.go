package controller

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"time"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	podv1alpha1 "github.com/zmotso/pod-calculator/api/v1alpha1"
)

// PodCalculatorReconciler reconciles a PodCalculator object
type PodCalculatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pod.example.com,resources=podcalculators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pod.example.com,resources=podcalculators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pod.example.com,resources=podcalculators/finalizers,verbs=update

//+kubebuilder:rbac:groups=v1,resources=configmap,verbs=get;list;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PodCalculatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)

	reqLogger.Info("Reconciling PodCalculator")

	calc := &podv1alpha1.PodCalculator{}
	err := r.Get(ctx, req.NamespacedName, calc)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			reqLogger.Info("PodCalculator resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	pods := corev1.PodList{}
	err = r.List(ctx, &pods, client.InNamespace(calc.Namespace), client.MatchingLabels(calc.Spec.LabelsSelector))
	if err != nil {
		return ctrl.Result{}, err
	}

	if calc.Spec.ConfigmapRef != nil {
		confM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      calc.Name,
				Namespace: calc.Namespace,
			},
		}

		if _, err = controllerutil.CreateOrUpdate(ctx, r.Client, confM, func() error {
			confM.Data = map[string]string{
				"pods": strconv.Itoa(len(pods.Items)),
			}

			return nil
		}); err != nil {
			return ctrl.Result{}, err
		}

		reqLogger.Info("ConfigMap with pods count created")
	}

	calc.Status.Count = len(pods.Items)
	if err = r.Status().Update(ctx, calc); err != nil {
		return ctrl.Result{}, err
	}

	reqLogger.Info("PodCalculator reconciled")

	return ctrl.Result{
		RequeueAfter: time.Second * 5,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodCalculatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podv1alpha1.PodCalculator{}).
		Complete(r)
}
