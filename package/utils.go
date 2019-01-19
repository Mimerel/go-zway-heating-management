package _package

import (
	"time"
)

func GetTimeAndDay(config *Configuration) {
	// getting time
	moment := time.Now()
	hour := moment.Hour() * 100
	config.Moment.Time = hour + moment.Minute()
	// getting weekday
	config.Moment.Weekday = moment.Weekday()
	// creatingDate
	config.Moment.Date = CreateDate(moment)
}

func CreateDate(moment time.Time) int {
	return moment.Year() * 10000 + int(moment.Month()) * 100 + moment.Day()
}