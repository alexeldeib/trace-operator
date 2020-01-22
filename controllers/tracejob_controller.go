package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/iovisor/kubectl-trace/pkg/meta"
	"github.com/iovisor/kubectl-trace/pkg/tracejob"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	observev1alpha1 "github.com/alexeldeib/trace-operator/api/v1alpha1"
)

var (
	// ImageNameTag represents the default tracerunner image
	DefaultImageNameTag = "quay.io/iovisor/kubectl-trace-bpftrace:latest"
	// InitImageNameTag represents the default init container image
	DefaultInitImageNameTag = "quay.io/iovisor/kubectl-trace-init:latest"
	// DefaultDeadline is the maximum time a tracejob is allowed to run, in seconds
	DefaultDeadline = 3600
	// DefaultDeadlineGracePeriod is the maximum time to wait to print a map or histogram, in seconds
	// note that it must account for startup time, as the deadline as based on start time
	DefaultDeadlineGracePeriod = 30
	// DefaultServiceAccount is the name of the default service account to run jobs with.
	DefaultServiceAccount = "default"
)

// TraceJobReconciler reconciles a TraceJob object
type TraceJobReconciler struct {
	client.Client
	RESTConfig *rest.Config
	Log        logr.Logger
	Scheme     *runtime.Scheme
}

type Object interface {
	metav1.Object
	runtime.Object
}

func build(in *observev1alpha1.TraceJob) tracejob.TraceJob {
	spec := tracejob.TraceJob{
		Name:                fmt.Sprintf("%s%s", meta.ObjectNamePrefix, in.Status.ID),
		Namespace:           in.Namespace,
		ServiceAccount:      DefaultServiceAccount,
		ID:                  in.Status.ID,
		Hostname:            in.Spec.Hostname,
		Program:             in.Spec.Program,
		ImageNameTag:        DefaultImageNameTag,
		InitImageNameTag:    DefaultInitImageNameTag,
		FetchHeaders:        in.Spec.FetchHeaders,
		Deadline:            int64(DefaultDeadline),
		DeadlineGracePeriod: int64(DefaultDeadlineGracePeriod),
	}

	if in.Spec.ServiceAccount != nil {
		spec.ServiceAccount = *in.Spec.ServiceAccount
	}
	if in.Spec.ImageNameTag != nil {
		spec.ImageNameTag = *in.Spec.ImageNameTag
	}
	if in.Spec.InitImageNameTag != nil {
		spec.InitImageNameTag = *in.Spec.InitImageNameTag
	}
	if in.Spec.Deadline != nil {
		spec.Deadline = *in.Spec.Deadline
	}
	if in.Spec.DeadlineGracePeriod != nil {
		spec.DeadlineGracePeriod = *in.Spec.DeadlineGracePeriod
	}

	return spec
}

// +kubebuilder:rbac:groups=observe.alexeldeib.xyz,resources=tracejobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=observe.alexeldeib.xyz,resources=tracejobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=jobs/status,verbs=get;update;patch

func (r *TraceJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("tracejob", req.NamespacedName)

	log.V(1).Info("fetching TraceJob")
	in := &observev1alpha1.TraceJob{}
	if err := r.Get(ctx, req.NamespacedName, in); err != nil {
		log.Error(err, "unable to fetch TraceJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Set trace id if not set
	if in.Status.ID == "" {
		juid := uuid.NewUUID()
		in.Status.ID = juid
		if err := r.Status().Update(ctx, in); err != nil {
			return ctrl.Result{}, err
		}
	}

	spec := build(in)

	job, cm := (&tracejob.TraceJobClient{}).CreateJob(spec)

	// Set owner ref on config map
	log.V(1).Info("setting owner on config map")
	if err := controllerutil.SetControllerReference(in, cm, r.Scheme); err != nil {
		log.Error(err, "failed to set controller owner reference on config map")
		return ctrl.Result{}, err
	}

	// Apply config map
	log.V(1).Info("applying config map")
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		return nil
	}); err != nil {
		log.Error(err, "failed to apply config map")
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(in, job, r.Scheme); err != nil {
		log.Error(err, "failed to set controller owner reference on job")
		return ctrl.Result{}, err
	}

	log.V(1).Info("applying job")
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, job, func() error {
		return nil
	}); err != nil {
		log.Error(err, "failed to apply job")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *TraceJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&observev1alpha1.TraceJob{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
