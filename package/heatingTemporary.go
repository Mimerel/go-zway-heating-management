package _package

import (
	"fmt"
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
		UpdateYamFile(config)
	} else {
		return fmt.Errorf("Wrong number of parameters sent")
	}
	return nil
}


