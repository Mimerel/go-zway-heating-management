package _package

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"net/http"
	"time"
)


func ExecuteRequest(config *Configuration, url string, id string, instance string, commandClass string, level float64) (err error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	postingUrl := "http://" + url + ":8083/ZWaveAPI/Run/devices[" + id + "].instances[" + instance + "].commandClasses[" + commandClass + "].Set(" + fmt.Sprintf("%f", level) + ")"
	logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Request posted : %s", postingUrl))


	_, err = client.Get(postingUrl)
	if err != nil {
		logs.Error(config.Elasticsearch.Url, config.Host, fmt.Sprintf("Failed to execute request %s \n", postingUrl, err))
		return err
	}
	return nil
}

func sendCommandToUpdateHeating(config *Configuration, value float64) (error) {
	url := config.GlobalSettings.ActualHeater.Ip
	id := config.GlobalSettings.ActualHeater.Id
	instance := config.GlobalSettings.ActualHeater.Instance
	commandClass := config.GlobalSettings.ActualHeater.CommandClass
	if url != "" {
		err := ExecuteRequest(config, url, id, instance, commandClass, value)
		if err != nil {
			return err
		}
	}
	return nil
}
