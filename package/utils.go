package _package

import (
	"time"
)

func GetTimeAndDay(config *Configuration) {
	// getting time
	config.Moment.Moment = time.Now()
	hour := config.Moment.Moment.Hour() * 100

	config.Moment.Time = hour + config.Moment.Moment.Minute()
	// getting weekday
	config.Moment.Weekday = config.Moment.Moment.Weekday()
	// creatingDate
	config.Moment.Date = CreateDate(config.Moment.Moment)
}

func CreateDate(moment time.Time) int {
	return moment.Year() * 10000 + int(moment.Month()) * 100 + moment.Day()
}

func StringInArray(value string, list []string) (bool) {
	for _,v := range list {
		if value == v {
			return true
		}
	}
	return false
}