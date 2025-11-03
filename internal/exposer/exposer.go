package exposer

import (
	"fmt"
	"itspeetah/np-dag-expo/internal/controller/neptuneplus"
	"net/http"

	"k8s.io/klog/v2"
)

type Exposer struct {
	httpConfig ExposerHttpConfig
	depDagCtrl *neptuneplus.DependencyGraphReconciler
}

type ExposerHttpConfig struct {
	Port int
	Host string
}

func NewExposer(
	httpConfig ExposerHttpConfig,
	ddr *neptuneplus.DependencyGraphReconciler,
) Exposer {
	return Exposer{
		httpConfig: httpConfig,
		depDagCtrl: ddr,
	}
}

func (e *Exposer) Run() error {

	klog.Info("Starting exposer...")
	go e.StartHttpServer()
	klog.Info("Started exposer.")
	return nil
}

func (e *Exposer) StartHttpServer() {
	config := e.httpConfig
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	klog.Infof("Staring http server on %s.", address)

	setupHttpRoutes(e)

	if err := http.ListenAndServe(address, nil); err != nil {
		klog.Error(err)
	}
}
