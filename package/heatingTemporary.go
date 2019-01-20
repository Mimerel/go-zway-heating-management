package _package

import (
	"fmt"
	"github.com/Mimerel/go-logger-client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func SettingTemporaryValues( config *Configuration, urlPath string) (err error) {
	GetTimeAndDay(config)
	urlParams := strings.Split(urlPath, "/")
	if len(urlParams) == 3 && strings.ToLower(urlParams[2])=="reset"  {
		config.TemporaryValues = Moment{}
	} else if len(urlParams) == 4 {
		hours, err := strconv.ParseInt(urlParams[3], 10, 64)
		if err != nil {
			return fmt.Errorf("unable to convert duration string to int64")
		}
		if !StringInArray(urlParams[2], config.GlobalSettings.AuthorizedLevels) {
			return fmt.Errorf("Level requested does not exist %s", urlParams[2])
		}
		config.TemporaryValues.Moment = config.Moment.Moment.Local().Add(time.Hour * time.Duration(hours))
		config.TemporaryValues.Level = urlParams[2]
		updateYamFile(config)
	} else {
		return fmt.Errorf("Wrong number of parameters sent")
	}
	return nil
}


func updateYamFile(config *Configuration)  {
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