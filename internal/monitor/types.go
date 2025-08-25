package monitor

type PodScaleData struct {
	Name             string
	Namespace        string
	Pod              string
	Service          string
	DesiredResources int64 //v1.ResourceList
	ActualResources  int64 //v1.ResourceList
	CappedResources  int64 //v1.ResourceList
}
