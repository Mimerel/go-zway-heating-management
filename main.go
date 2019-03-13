package main

import (
	"encoding/json"
	"go-zway-heating-management/package"
	"html/template"
	"net/http"
)

func main() {

	config := _package.ReadConfiguration()
	Port := config.Port

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		err := _package.UpdateHeating(w, r, &config)
		if err != nil {
			config.Logger.Error("Unable to update heating %+v ", err)
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
		err := _package.SettingTemporaryValues(&config, r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			t := template.New("confirmation.html")
			t, err := t.ParseFiles(config.GlobalSettings.ApplicationRunningPath + "/package/templates/confirmation.html")
			if err != nil {
				config.Logger.Error("Error Parsing template%+v", err)
			}
			err = t.Execute(w, _package.Confirmation{
				IpPort: config.Ip + ":" + config.Port,
			} )
			if err != nil {
				config.Logger.Error("Error Execution %+v", err)
			}
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
		config.Logger.Error("Error Parsing template%+v", err)
	}
	data, err := _package.HeatingStatus(config)
	err = t.Execute(w, data)
	if err != nil {
		config.Logger.Error("Error Execution %+v", err)
	}
}
