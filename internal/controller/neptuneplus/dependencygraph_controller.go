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

package neptuneplus

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/modern-go/concurrent"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	neptuneplusv1alpha1 "itspeetah/np-tester/api/neptuneplus/v1alpha1"
)

// DependencyGraphReconciler reconciles a DependencyGraph object
type DependencyGraphReconciler struct {
	client.Client
	Scheme                *runtime.Scheme
	ExternalResponseTimes *concurrent.Map
}

// +kubebuilder:rbac:groups=neptuneplus.polimi.it,resources=dependencygraphs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=neptuneplus.polimi.it,resources=dependencygraphs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=neptuneplus.polimi.it,resources=dependencygraphs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DependencyGraph object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *DependencyGraphReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	klog.Infof("[DependencyGraph] Reconcile loop triggered by %s", req.NamespacedName)
	depDag := &neptuneplusv1alpha1.DependencyGraph{}
	err := r.Get(ctx, req.NamespacedName, depDag)
	if err != nil {
		if apierrors.IsNotFound(err) {
			klog.Info("dependencygraph resource not found. Ignoring since object must be deleted")
			r.ExternalResponseTimes.Delete(req.NamespacedName)
			return ctrl.Result{}, nil
		}
		klog.Error(err, "Failed to get dependencygraph")
		return ctrl.Result{}, err
	}

	for _, node := range depDag.Status.Nodes {
		key := fmt.Sprintf("%s/%s", node.FunctionNamespace, node.FunctionName)
		r.ExternalResponseTimes.Store(key, node.ExternalResponseTime)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DependencyGraphReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&neptuneplusv1alpha1.DependencyGraph{}).
		Named("neptuneplus-dependencygraph").
		Complete(r)
}
