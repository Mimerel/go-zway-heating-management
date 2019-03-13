package _package

func HeatingStatus(config *Configuration) (data Status, err error) {
	GetTimeAndDay(config)

	floatLevel, err := GetInitialHeaterParams(config)
	if err != nil {
		floatLevel = 15
	}
	heater, temperature := collectMetrics(config)

	data.Until = config.TemporaryValues.Moment
	data.Temperature_Actual = temperature
	data.Temperature_Requested = floatLevel
	data.Heater_Level = heater
	data.TemporaryLevel = config.TemporaryValues.Level
	if config.TemporaryValues.Level != "" {
		data.IsTemporary = true
	} else {
		data.IsTemporary = false
	}
	data.IpPort = config.Ip + ":" + config.Port
	data.UpdateTime = config.GlobalSettings.LastUpdate
	data.NormalValues = config.NormalValues
	return data, nil
}

func collectMetrics(config *Configuration) (heater float64, temperature float64) {
	err := getAllActualMetricValues(config)
	if err != nil {
		config.Logger.Error("Unable to get actual metric values", err)
		return  0, 9999.0
	} else {
		err = getRequiredMetrics(config)
		if err != nil {
			return  0, 9999.0
			config.Logger.Error("Cound not find actual metrics required", err)
		} else {
			heater = config.GlobalSettings.ActualHeater.Value
			temperature = config.GlobalSettings.ActualTemperature.Value
		}
	}
	config.Logger.Info("Metrics retrieved, heater %f , temperature %f", heater, temperature)
	return heater, temperature
}
