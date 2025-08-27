package exposer

import (
	"bytes"
	"fmt"
	"itspeetah/np-tester/internal/monitor"
	"net/http"

	"k8s.io/klog/v2"
)

func setupHttpRoutes(e *Exposer) {
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello, world!\n")
	})

	http.HandleFunc("/headers", func(w http.ResponseWriter, req *http.Request) {
		for name, headers := range req.Header {
			for _, h := range headers {
				fmt.Fprintf(w, "%v: %v\n", name, h)
			}
		}
	})

	http.HandleFunc("/latest-log", func(w http.ResponseWriter, req *http.Request) {
		klog.Info("[EXPOSER] requested latest log")
		txt := e.getLatestLog()
		fmt.Fprint(w, txt)
	})

	http.HandleFunc("/reset-cached", func(w http.ResponseWriter, req *http.Request) {
		klog.Info("[EXPOSER] requested cached value reset")
		e.clearCachedValues()
		fmt.Fprint(w, "OK")
	})
}

func (e *Exposer) getLatestLog() string {
	var buffer bytes.Buffer

	buffer.WriteString("PODSCALES\n")
	e.podscaleData.Range(func(key any, value any) bool {

		keyAsStr, _ := key.(string)
		podScaleData, _ := value.(monitor.PodScaleData)

		buffer.WriteString(
			fmt.Sprintf(
				"%s %s %s %s %d %d %d\n",
				keyAsStr,
				podScaleData.Name,
				podScaleData.Pod,
				podScaleData.Service,
				podScaleData.DesiredResources,
				podScaleData.ActualResources,
				podScaleData.CappedResources,
			))

		return true
	})
	buffer.WriteString("EXT. RES. TIMES\n")
	e.depDagExtRTs.Range(func(key any, value any) bool {

		keyAsStr, _ := key.(string)
		valueAsInt64, _ := value.(int64)

		buffer.WriteString(fmt.Sprintf("%s %d\n", keyAsStr, valueAsInt64))

		return true
	})
	buffer.WriteString("POD RES. TIMES\n")
	e.podResponseTimes.Range(func(key any, value any) bool {
		keyAsStr, _ := key.(string)
		valueAsInt64, _ := value.(int64)

		buffer.WriteString(fmt.Sprintf("%s %d\n", keyAsStr, valueAsInt64))

		return true
	})

	return buffer.String()
}

func (e *Exposer) clearCachedValues() {
	e.depDagExtRTs.Clear()
	e.podResponseTimes.Clear()
	e.podscaleData.Clear()
}
