package main

import (
	"encoding/json"
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"github.com/op/go-logging"
	"go-zway-heating-management/package"
	"html/template"
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
			_package.UpdateYamFile(&config)
		}
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		StatusPage(w, r, &config)
	})

	http.HandleFunc("/temporary/", func(w http.ResponseWriter, r *http.Request) {
		err :=	_package.SettingTemporaryValues(&config, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			StatusPage(w, r, &config)
		}
	})


	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		data, err := _package.HeatingStatus(&config)
		var js []byte
		js, err = json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})



	http.ListenAndServe(":"+Port, nil)
}

func StatusPage(w http.ResponseWriter, r *http.Request, config *_package.Configuration) {
	t := template.New("status.html")
	t, err := t.ParseFiles(config.GlobalSettings.ApplicationRunningPath + "/package/templates/status.html")
	if err != nil {
		logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Error Parsing template%+v", err))
	}
	data, err := _package.HeatingStatus(config)
	err = t.Execute(w, data)
	if err != nil {
		logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Error Execution %+v", err))
	}
}
