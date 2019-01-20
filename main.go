package main

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"github.com/op/go-logging"
	"go-zway-heating-management/package"
	"net/http"
	"os"
)

var log = logging.MustGetLogger("default")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
)

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(backendLeveled, backendFormatter)

	config := _package.ReadConfiguration()
	Port := config.Port
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		err := _package.UpdateHeating(w, r, &config)
		if err != nil {
			logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Unable to update heating %+v ", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(200)
		}
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		err := _package.HeatingStatus(w, r, &config)
		if err != nil {
			logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Unable to update heating %+v ", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(200)
		}
	})
	http.ListenAndServe(":"+Port, nil)
}


