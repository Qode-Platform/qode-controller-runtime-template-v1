// Package controller holds the reconcilers.
package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ConfigMapReconciler stamps a label on every ConfigMap it sees. Replace the
// body with the real reconciliation; the shape is the point.
type ConfigMapReconciler struct {
	Client client.Client
}

const ManagedLabel = "qode.world/managed-by"

func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	var cm corev1.ConfigMap
	if err := r.Client.Get(ctx, req.NamespacedName, &cm); err != nil {
		// Gone is not an error: the object was deleted between the event and now.
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	if cm.Labels[ManagedLabel] == "qode" {
		return reconcile.Result{}, nil // already correct, do nothing
	}
	if cm.Labels == nil {
		cm.Labels = map[string]string{}
	}
	cm.Labels[ManagedLabel] = "qode"
	log.FromContext(ctx).Info("labelling configmap", "name", req.NamespacedName)
	return reconcile.Result{}, r.Client.Update(ctx, &cm)
}
