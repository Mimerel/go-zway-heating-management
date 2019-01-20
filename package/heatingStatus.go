package _package

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
)

func HeatingStatus( config *Configuration) (data Status, err error) {
	GetTimeAndDay(config)

	floatLevel, heater, temperature, err := GetInitialHeaterParams(config)

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

	data.Until = config.TemporaryValues.Moment
	data.Temperature_Actual = temperature
	data.Temperature_Requested = floatLevel
	data.Heater_Level = heater
	data.TemporaryLevel = config.TemporaryValues.Level
	if 	config.TemporaryValues.Level != "" {
		data.IsTemporary = true
	}	else {
		data.IsTemporary = false
	}
	data.IpPort = config.Ip + ":" + config.Port
	data.UpdateTime = config.Moment.Moment
	return data, nil
}

