package main

import (
	"encoding/json"
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"github.com/op/go-logging"
	"go-zway-heating-management/package"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
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
		err := updateHeating(w, r, &config)
		if err != nil {
			logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Unable to update heating %+v ", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(200)
		}
	})
	http.ListenAndServe(":"+Port, nil)
}

func updateHeating(w http.ResponseWriter, r *http.Request, config *_package.Configuration) (error) {
	_package.GetTimeAndDay(config)
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
			heater = config.ActualHeater.Value
			temperature = config.ActualTemperature.Value
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

func ExecuteRequest(url string, id string, instance string, commandClass string, level string) (err error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + id + "].instances[" + instance + "].commandClasses[" + commandClass + "].Set(" + level + ")"
	log.Info("Request posted : %s", postingUrl)

	_, err = client.Get(postingUrl)
	if err != nil {
		fmt.Printf("Failed to execute request %s \n", postingUrl, err)
		return err
	}
	return nil
}

func sendCommandToUpdateHeating(config *_package.Configuration, value float64) (error) {
	url := config.ActualHeater.Ip
	id := config.ActualHeater.Id
	instance := config.ActualHeater.Instance
	commandClass := config.ActualHeater.CommandClass
	if url != "" {
		err := ExecuteRequest(url, id, instance, commandClass, strconv.FormatFloat(value, 'g', 1, 64))
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckIfHeatingNeedActivating(config *_package.Configuration, floatLevel float64, temperature float64) bool {
	if temperature < floatLevel {
		return true
	}
	return false
}

func getRequiredMetrics(config *_package.Configuration) (error) {
	found := 0
	for _, v := range config.Metrics {
		if v.Metric == config.GlobalSettings.HeaterMetricName {
			value, err := strconv.ParseFloat(v.Value, 64)
			if err == nil {
				config.ActualHeater.Value = value
				found += 1
			}
		}
		if v.Metric == config.GlobalSettings.TemperatureMetricName {
			value, err := strconv.ParseFloat(v.Value, 64)
			config.ActualTemperature.Value = value
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

func getValueOfLevel(config *_package.Configuration, setLevel string) (float64) {
	level := 15.0
	for _, v := range config.GlobalSettings.Levels {
		if v.Name == setLevel {
			return v.Level
		}
	}
	return level
}

func getLevel(config *_package.Configuration, ) (string, error) {
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

func getLevelValue(config *_package.Configuration, index int) (string, error) {
	for _, v := range config.NormalValues[index].Settings {
		if v.From < config.Moment.Time && v.To >= config.Moment.Time {
			return v.Level, nil
		}
	}
	return "away", fmt.Errorf("Unable to find normal level for weekday %s, time %d", config.Moment.Weekday, config.Moment.Time)
}

func getAllActualMetricValues(config *_package.Configuration) (error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := config.LastMetricsEndpoint.Url + "/all"
	response, err := client.Get(postingUrl)
	if err != nil {
		fmt.Printf("Failed to retrieve metrics")
		return fmt.Errorf("Failed to retrieve metrics", err)
	}
	var metrics []_package.StructuredData

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
