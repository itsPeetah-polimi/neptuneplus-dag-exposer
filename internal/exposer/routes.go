package exposer

import (
	"fmt"
	"net/http"
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

	http.HandleFunc("/get-json", e.depDagCtrl.HttpGETGraphJson)
}
