/*
Copyright 2025.

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

package systemautoscaler

import (
	"context"
	"fmt"

	"github.com/modern-go/concurrent"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	systemautoscalerv1beta1 "itspeetah/np-tester/api/systemautoscaler/v1beta1"
	"itspeetah/np-tester/internal/monitor"
)

// PodScaleReconciler reconciles a PodScale object
type PodScaleReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	PodScaleData *concurrent.Map
}

// +kubebuilder:rbac:groups=systemautoscaler.polimi.it,resources=podscales,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=systemautoscaler.polimi.it,resources=podscales/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=systemautoscaler.polimi.it,resources=podscales/finalizers,verbs=update

// +kubebuilder:rbac:groups=custom.metrics.k8s.io,resources=pods/response_time,verbs=get
// +kubebuilder:rbac:groups=core,resources=services;pods,verbs=get;list;watch;

func (r *PodScaleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)
	klog.Infof("[PodScale] Reconcile loop triggered by %s", req.NamespacedName)

	podScale := &systemautoscalerv1beta1.PodScale{}
	err := r.Get(ctx, req.NamespacedName, podScale)
	if err != nil {
		if apierrors.IsNotFound(err) {
			klog.Info("podscale resource not found. Ignoring since object must be deleted")
			r.PodScaleData.Delete(req.NamespacedName)
			return ctrl.Result{}, nil
		}
		klog.Error(err)
		return ctrl.Result{}, err
	}

	podScaleData := &monitor.PodScaleData{
		Name:             podScale.Name,
		Namespace:        podScale.Namespace,
		Pod:              podScale.Spec.Pod,
		Service:          podScale.Spec.Service,
		DesiredResources: podScale.Spec.DesiredResources.Cpu().MilliValue(),
		ActualResources:  podScale.Status.ActualResources.Cpu().MilliValue(),
		CappedResources:  podScale.Status.CappedResources.Cpu().MilliValue(),
	}

	key := fmt.Sprintf("%s/%s", podScale.Namespace, podScale.Name)
	r.PodScaleData.Store(key, *podScaleData)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodScaleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&systemautoscalerv1beta1.PodScale{}).
		Named("podscale").
		Complete(r)
}
