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
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"setcreed.github.io/store/internal/builders"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "setcreed.github.io/store/api/v1alpha1"
)

const (
	Kind            = "DbConfig"
	GroupAPIVersion = "apps.setcreed.github.io/v1alpha1"
)

// DbConfigReconciler reconciles a DbConfig object
type DbConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.setcreed.github.io,resources=dbconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.setcreed.github.io,resources=dbconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.setcreed.github.io,resources=dbconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DbConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *DbConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	config := &appsv1alpha1.DbConfig{}
	err := r.Get(ctx, req.NamespacedName, config)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	builder, err := builders.NewDeployBuilder(config, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = builder.Build(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *DbConfigReconciler) OnDelete(ctx context.Context, event event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.Object.GetOwnerReferences() {
		if ref.Kind == Kind && ref.APIVersion == GroupAPIVersion {
			limitingInterface.AddRateLimited(reconcile.Request{types.NamespacedName{
				Namespace: event.Object.GetNamespace(),
				Name:      ref.Name,
			}})
		}
	}
}

func (r *DbConfigReconciler) OnUpdate(ctx context.Context, event event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.ObjectNew.GetOwnerReferences() {
		if ref.Kind == Kind && ref.APIVersion == GroupAPIVersion {
			limitingInterface.AddRateLimited(reconcile.Request{types.NamespacedName{
				Namespace: event.ObjectNew.GetNamespace(),
				Name:      ref.Name,
			}})
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DbConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.DbConfig{}).
		Watches(&appv1.Deployment{}, handler.Funcs{
			DeleteFunc: r.OnDelete,
			UpdateFunc: r.OnUpdate,
		}).
		Complete(r)
}
