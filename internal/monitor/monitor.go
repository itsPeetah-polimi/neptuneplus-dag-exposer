package monitor

import (
	"time"

	metrics "itspeetah/np-tester/internal/metrics"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/modern-go/concurrent"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
)

type Monitor struct {
	Client           client.Client
	MetricClient     metrics.MetricGetter
	podscaleData     *concurrent.Map
	depDagExtRTs     *concurrent.Map
	podResponseTimes *concurrent.Map
}

func NewMonitor(
	client client.Client,
	metricClient metrics.MetricGetter,
	inPodscales *concurrent.Map,
	extRTMap *concurrent.Map,
	podResponseTimes *concurrent.Map,
) Monitor {
	return Monitor{
		Client:           client,
		MetricClient:     metricClient,
		podscaleData:     inPodscales,
		depDagExtRTs:     extRTMap,
		podResponseTimes: podResponseTimes,
	}
}

func (m *Monitor) Run() error {

	klog.Info("Starting monitor controller...")
	go wait.Forever(m.runPodScaleLoop, 4*time.Second)
	go wait.Forever(m.runDependencyGraphLoop, 4*time.Second)
	go wait.Forever(m.runMetricGetterLoop, 1*time.Second)
	klog.Info("Started monitor controller.")
	return nil
}

func (m *Monitor) runPodScaleLoop() {
	klog.Info("PodScale Monitor Loop")
	m.podscaleData.Range(func(key any, value any) bool {
		// keyAsStr, _ := key.(string)
		podScaleData, _ := value.(PodScaleData)

		klog.Infof("[PodScl] pod %s (%s) cpu (d: %d, a: %d, c: %d)",
			podScaleData.Pod,
			podScaleData.Service,
			podScaleData.DesiredResources,
			podScaleData.ActualResources,
			podScaleData.CappedResources,
		)

		return true
	})
}

func (m *Monitor) runDependencyGraphLoop() {
	klog.Info("DependencyGraph Monitor Loop")
	m.depDagExtRTs.Range(func(key any, value any) bool {

		keyAsStr, _ := key.(string)
		valueAsInt64, _ := value.(int64)

		klog.Infof("[DepDag] ert for %s:\t\t%d", keyAsStr, valueAsInt64)

		return true
	})
}

func (m *Monitor) runMetricGetterLoop() {

	m.podscaleData.Range(func(key any, value any) bool {
		// keyAsStr, _ := key.(string)
		podScaleData, _ := value.(PodScaleData)
		podName := podScaleData.Pod

		podRt := m.getPodResponseTime(podName)
		m.podResponseTimes.Store(podName, podRt)

		klog.Infof("[PodRts] rt for %s: %d", podName, podRt)

		return true
	})
}
