package _package

import "time"

type Configuration struct {
	Token string `yaml:"token,omitempty"`
	Elasticsearch URLS `yaml:"elasticSearch,omitempty"`
	LastMetricsEndpoint URLS `yaml:"metricsEnpoint,omitempty"`
	Host string `yaml:"host,omitempty"`
	Port string `yaml:"port,omitempty"`
	GlobalSettings GlobalSettingsType `yaml:"settings,omitempty"`
	NormalValues []Normal `yaml:"normal,omitempty"`
	Metrics []StructuredData
	ActualHeater ZwaveParams
	ActualTemperature ZwaveParams
	Moment Moment
}

type ZwaveParams struct {
	Name string `yaml:"name,omitempty"`
	Ip string `yaml:"ip,omitempty"`
	Id string `yaml:"id,omitempty"`
	CommandClass string `yaml:"commandClass,omitempty"`
	Instance string `yaml:"instance,omitempty"`
	Value float64
}

type URLS struct {
	Url string `yaml:"url,omitempty"`
}

type StructuredData struct {
	Metric string
	Labels map[string]string
	Timestamp string
	Timestamp2 string
	Value string
}

type GlobalSettingsType struct {
	Levels []Level `yaml:"level,omitempty"`
	HeaterMetricName string `yaml:"heaterMetricName,omitempty"`
	TemperatureMetricName string `yaml:"temperatureMetricName,omitempty"`
}

type Level struct {
	Name string `yaml:"name,omitempty"`
	Level float64 `yaml:"level,omitempty"`
}


type Normal struct {
	Weekday time.Weekday `yaml:"day,omitempty"`
	Name string `yaml:"name,omitempty"`
	Settings []Period `yaml:"settings,omitempty"`
}


type Period struct {
	From int `yaml:"from,omitempty"`
	To int `yaml:"to,omitempty"`
	Level string `yaml:"level,omitempty"`
}

type Moment struct {
	Time int
	Weekday time.Weekday
	Date int
}