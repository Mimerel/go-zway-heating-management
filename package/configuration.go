package _package

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

/**
Reads configuration file
 */
func ReadConfiguration() (Configuration) {
	pathToFile := os.Getenv("LOGGER_CONFIGURATION_FILE")
	if _, err := os.Stat("./configuration.yaml"); !os.IsNotExist(err) {
		pathToFile = "./configuration.yaml"
	} else if pathToFile == "" {
		pathToFile = "/home/pi/go/src/go-zway-heating-management/configuration.yaml"
	}
	yamlFile, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		panic(err)
	}

	var config Configuration

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	} else {
		config.GlobalSettings.ApplicationRunningPath = strings.Replace(pathToFile, "configuration.yaml", "", 1)
		config.GlobalSettings.AuthorizedLevels = getAllSetLevels(config.GlobalSettings.Levels)
		fmt.Printf("Configuration Loaded : %+v \n", config)
	}
	return config
}

func getAllSetLevels(levels []Level) (result []string) {
	for _, v := range levels {
		result = append(result, v.Name)
	}
	return result
}
