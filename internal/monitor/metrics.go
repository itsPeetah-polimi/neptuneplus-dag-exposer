package monitor

import (
	"context"
	"itspeetah/np-tester/internal/metrics"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (m *Monitor) getPodResponseTime(podName string) int64 {
	pod := &v1.Pod{}
	key := client.ObjectKey{
		Namespace: "openfaas-fn",
		Name:      podName,
	}

	if err := m.Client.Get(context.TODO(), key, pod); err != nil {
		klog.Error(err)
		return 0
	}

	rt := metrics.PodResponseTime(m.MetricClient, pod)
	return rt
}
