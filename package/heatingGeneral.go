package _package

import (
	"encoding/json"
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func getLevel(config *Configuration, ) (string, error) {
	setLevel := "away"
	for key, v := range config.NormalValues {
		if v.Weekday == config.Moment.Weekday {
			foundValue, err := getLevelValue(config, key)
			if err != nil {
				logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Unable to find temperature ", err))
			} else {
				setLevel = foundValue
			}
		}
	}
	return setLevel, nil
}

func getLevelValue(config *Configuration, index int) (string, error) {
	for _, v := range config.NormalValues[index].Settings {
		if v.From < config.Moment.Time && v.To >= config.Moment.Time {
			return v.Level, nil
		}
	}
	return "away", fmt.Errorf("Unable to find normal level for weekday %s, time %d", config.Moment.Weekday, config.Moment.Time)
}

func getAllActualMetricValues(config *Configuration) (error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := config.LastMetricsEndpoint.Url + "/all"
	response, err := client.Get(postingUrl)
	if err != nil {
		return fmt.Errorf("Failed to retrieve metrics", err)
	}
	var metrics []StructuredData

	defer response.Body.Close()

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read body of metrics request", err)
	}

	err = json.Unmarshal(buf, &metrics)
	if err != nil {
		return fmt.Errorf("Failed to convert body to json", err)
	}

	config.Metrics = metrics

	return nil
}


func CheckIfHeatingNeedsActivating(config *Configuration, floatLevel float64, temperature float64) bool {
	if temperature < floatLevel {
		return true
	}
	return false
}

func getRequiredMetrics(config *Configuration) (error) {
	found := 0
	for _, v := range config.Metrics {
		if v.Metric == config.GlobalSettings.ActualHeater.Name {
			value, err := strconv.ParseFloat(v.Value, 64)
			if err == nil {
				config.GlobalSettings.ActualHeater.Value = value
				found += 1
			}
		}
		if v.Metric == config.GlobalSettings.ActualTemperature.Name {
			value, err := strconv.ParseFloat(v.Value, 64)
			config.GlobalSettings.ActualTemperature.Value = value
			if err == nil {
				found += 1
			}
		}
	}
	if found != 2 {
		return fmt.Errorf("Only found actual values for %d out of 2 required", found)
	}
	return nil
}

func getValueOfLevel(config *Configuration, setLevel string) (float64) {
	level := 15.0
	for _, v := range config.GlobalSettings.Levels {
		if v.Name == setLevel {
			return v.Level
		}
	}
	return level
}

func GetInitialHeaterParams(config *Configuration) (floatLevel float64, err error) {
	setLevel, err := getLevel(config)
	logs.Info(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Retreived heating level, %v", setLevel))

	if err != nil {
		return floatLevel,err
	}
	if config.TemporaryValues.Moment.After(config.Moment.Moment) {
		setLevel = config.TemporaryValues.Level
		logs.Info(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Temporary heating override, %v", setLevel))
	} else if config.TemporaryValues.Moment.Before(config.Moment.Moment) {
		config.TemporaryValues = Moment{}
		logs.Info(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Clearing old temporary settings"))
	}
	return getValueOfLevel(config, setLevel), nil
}

func UpdateYamFile(config *Configuration)  {
	yamlFile, err := yaml.Marshal(config)
	if err != nil {
		logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Unable to yaml marshal local_storage file %+v", err))
	}
	err = ioutil.WriteFile(config.GlobalSettings.ApplicationRunningPath + "configuration.yaml", yamlFile, 0777)
	if err != nil {
		logs.Error("", "", fmt.Sprintf("Unable to write configuration file %+v", err))
	} else {
		logs.Info(config.Elasticsearch.Url, config.Host, "Configuration file updated\n")
	}
}