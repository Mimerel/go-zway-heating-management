package _package

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"net/http"
)

func UpdateHeating(w http.ResponseWriter, r *http.Request, config *Configuration) (error) {
	GetTimeAndDay(config)
	config.GlobalSettings.LastUpdate = config.Moment.Moment
	floatLevel, err := GetInitialHeaterParams(config)
	if err != nil {
		floatLevel = 15
	}
	heater, temperature := collectMetrics(config)

	activateHeating := CheckIfHeatingNeedsActivating(config, floatLevel, temperature)
	logs.Info(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Heating should be activated, %t", activateHeating))
	if heater == 0 && activateHeating {
		err = sendCommandToUpdateHeating(config, 255)
		if err != nil {
			return err
		}
	}
	if heater == 255 && !activateHeating {
		err = sendCommandToUpdateHeating(config, 0)
		if err != nil {
			return err
		}
	}
	return nil
}
