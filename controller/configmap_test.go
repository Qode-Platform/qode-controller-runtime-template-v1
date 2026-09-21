package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileLabelsConfigMap(t *testing.T) {
	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "demo", Namespace: "default"}}
	c := fake.NewClientBuilder().WithObjects(cm).Build()
	r := &ConfigMapReconciler{Client: c}

	key := types.NamespacedName{Name: "demo", Namespace: "default"}
	if _, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: key}); err != nil {
		t.Fatal(err)
	}

	var got corev1.ConfigMap
	if err := c.Get(context.Background(), key, &got); err != nil {
		t.Fatal(err)
	}
	if got.Labels[ManagedLabel] != "qode" {
		t.Fatalf("label = %q, want qode", got.Labels[ManagedLabel])
	}
}

func TestReconcileMissingObjectIsNotAnError(t *testing.T) {
	c := fake.NewClientBuilder().Build()
	r := &ConfigMapReconciler{Client: c}
	key := types.NamespacedName{Name: "gone", Namespace: "default"}
	if _, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: key}); err != nil {
		t.Fatalf("deleted object should not error: %v", err)
	}
}
