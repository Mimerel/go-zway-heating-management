package _package

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"net/http"
)

func UpdateHeating(w http.ResponseWriter, r *http.Request, config *Configuration) (error) {
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
			logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Cound not found actual metrics required", err))
		} else {
			heater = config.GlobalSettings.ActualHeater.Value
			temperature = config.GlobalSettings.ActualTemperature.Value
		}
	}

	activateHeating := CheckIfHeatingNeedActivating(config, FloatLevel, temperature)
	if heater == 0 && activateHeating {
		err = sendCommandToUpdateHeating(config, 255)
	}
	if heater == 255 && !activateHeating {
		err = sendCommandToUpdateHeating(config, 0)
	}
	return nil
}
