package _package

import (
	"encoding/json"
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"net/http"
)

func HeatingStatus(w http.ResponseWriter, r *http.Request, config *Configuration) (error) {
	GetTimeAndDay(config)
	setLevel, err := getLevel(config)
	if err != nil {
		return err
	}
	FloatLevel := getValueOfLevel(config, setLevel)
	heater := 255.0
	temperature := 9999.0

	// Getting actual metrics and values for required metrics
	err = getAllActualMetricValues(config)
	if err != nil {
		logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Unable get actual metric values", err))
	} else {
		err = getRequiredMetrics(config)
		if err != nil {
			logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Cound not find actual metrics required", err))
		} else {
			heater = config.GlobalSettings.ActualHeater.Value
			temperature = config.GlobalSettings.ActualTemperature.Value
		}
	}

	data := Status{	heater, FloatLevel, temperature}
	var js []byte
	js, err = json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}
