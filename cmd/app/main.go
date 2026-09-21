// A controller-runtime manager.
//
// NOT RUNNABLE OUTSIDE A CLUSTER: the manager needs a kubeconfig or an
// in-cluster service account, so START_CMD is empty and bin/run stops at the
// start step. Under k8s the manager does serve /healthz and /readyz on
// $PORT — wire HEALTH_PATH back up if you deploy it as a workload.
package main

import (
	"log"
	"os"
	"strings"

	"github.com/Qode-Platform/qode-controller-runtime-template-v1/controller"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

func port() string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return p
	}
	return "8081"
}

func main() {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		HealthProbeBindAddress: ":" + port(),
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Fatal(err)
	}
	if err := ctrl.NewControllerManagedBy(mgr).For(&corev1.ConfigMap{}).
		Complete(&controller.ConfigMapReconciler{Client: mgr.GetClient()}); err != nil {
		log.Fatal(err)
	}
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal(err)
	}
}
