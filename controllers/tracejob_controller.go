package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/iovisor/kubectl-trace/pkg/meta"
	"github.com/iovisor/kubectl-trace/pkg/tracejob"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	observev1alpha1 "github.com/alexeldeib/trace-operator/api/v1alpha1"
)

// TraceJobReconciler reconciles a TraceJob object
type TraceJobReconciler struct {
	client.Client
	RESTConfig *rest.Config
	Log        logr.Logger
	Scheme     *runtime.Scheme
}

// +kubebuilder:rbac:groups=observe.alexeldeib.xyz,resources=tracejobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=observe.alexeldeib.xyz,resources=tracejobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=jobs/status,verbs=get;update;patch

func (r *TraceJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("tracejob", req.NamespacedName)

	in := &observev1alpha1.TraceJob{}
	if err := r.Get(ctx, req.NamespacedName, in); client.IgnoreNotFound(err) != nil {
		log.Error(err, "unable to fetch TraceJob")
		return ctrl.Result{}, err
	}

	if in.Status.ID == "" {
		in.Status.ID = uuid.NewUUID()
		if err := r.Update(ctx, in); err != nil {
			log.Error(err, "unable to set ID for TraceJob")
			return ctrl.Result{}, err
		}
	}

	spec := tracejob.TraceJob{
		Name:                fmt.Sprintf("%s%s", meta.ObjectNamePrefix, string(in.Status.ID)),
		Namespace:           in.Namespace,
		ServiceAccount:      in.Spec.ServiceAccount,
		ID:                  in.Status.ID,
		Hostname:            in.Spec.Hostname,
		Program:             in.Spec.Program,
		PodUID:              in.Spec.PodUID,
		ContainerName:       in.Spec.ContainerName,
		IsPod:               in.Spec.IsPod,
		ImageNameTag:        in.Spec.ImageNameTag,
		InitImageNameTag:    in.Spec.InitImageNameTag,
		FetchHeaders:        in.Spec.FetchHeaders,
		Deadline:            in.Spec.Deadline,
		DeadlineGracePeriod: in.Spec.DeadlineGracePeriod,
	}

	job, cm := (&tracejob.TraceJobClient{}).CreateJob(spec)
	objects := []runtime.Object{job, cm}
	for _, obj := range objects {
		if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, obj, func() error {
			return nil
		}); err != nil {
			log.Error(err, "failed to apply object", "gvk", obj.GetObjectKind().GroupVersionKind().String())
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *TraceJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&observev1alpha1.TraceJob{}).
		Complete(r)
}
