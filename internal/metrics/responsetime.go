package metrics

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

func PodResponseTime(client MetricGetter, pod *corev1.Pod) int64 {
	rt, err := client.PodMetrics(pod, ResponseTime)
	if err != nil {
		klog.Errorf("[%s] Could not retrieve pod metrics for pod: %s", pod.Name, err.Error())
		return 0
	}
	return rt.Value.MilliValue()
}
