package _package

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
)

func HeatingStatus( config *Configuration) (data Status, err error) {
	GetTimeAndDay(config)
	setLevel, err := getLevel(config)
	if err != nil {
		return data, err
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

	return Status{	Heater_Level: heater, Temperature_Requested:FloatLevel, Temperature_Actual: temperature}, nil
}
