package exposer

import (
	"fmt"
	"net/http"

	"github.com/modern-go/concurrent"
	"k8s.io/klog/v2"
)

type Exposer struct {
	podscaleData     *concurrent.Map
	depDagExtRTs     *concurrent.Map
	podResponseTimes *concurrent.Map
	httpConfig       ExposerHttpConfig
}

type ExposerHttpConfig struct {
	Port int
	Host string
}

func NewExposer(
	httpConfig ExposerHttpConfig,
	inPodscales *concurrent.Map,
	extRTMap *concurrent.Map,
	podResponseTimes *concurrent.Map,
) Exposer {
	return Exposer{
		httpConfig:       httpConfig,
		podscaleData:     inPodscales,
		depDagExtRTs:     extRTMap,
		podResponseTimes: podResponseTimes,
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
