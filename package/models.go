package _package

import "time"

type Configuration struct {
	Elasticsearch URLS `yaml:"elasticSearch"`
	LastMetricsEndpoint URLS `yaml:"metricsEnpoint,omitempty"`
	Host string `yaml:"host,omitempty"`
	Ip string `yaml:"ip,omitempty"`
	Port string `yaml:"port,omitempty"`
	GlobalSettings GlobalSettingsType `yaml:"settings,omitempty"`
	NormalValues []Normal `yaml:"normal,omitempty"`
	TemporaryValues Moment `yaml:"temporary,omitempty"`
	Metrics []StructuredData
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
	Levels []Level `yaml:"levels,omitempty"`
	ActualHeater ZwaveParams `yaml:"heaterMetric,omitempty"`
	ActualTemperature ZwaveParams `yaml:"temperatureMetric,omitempty"`
	ApplicationRunningPath string
	AuthorizedLevels []string
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
	From int `yaml:"from"`
	To int `yaml:"to"`
	Level string `yaml:"level,omitempty"`
}

type Moment struct {
	Moment time.Time
	Time int
	Weekday time.Weekday
	Date int
	Level string
}

type Status struct {
	Heater_Level float64
	Temperature_Requested float64
	Temperature_Actual float64
	Until time.Time
	TemporaryLevel string
	IsTemporary bool
	IpPort string
}