package _package

import "time"

type Configuration struct {
	Elasticsearch URLS `yaml:"elasticSearch"`
	LastMetricsEndpoint URLS `yaml:"metricsEnpoint"`
	Host string `yaml:"host,omitempty"`
	Ip string `yaml:"ip,omitempty"`
	Port string `yaml:"port,omitempty"`
	GlobalSettings GlobalSettingsType `yaml:"settings,omitempty"`
	NormalValues []Normal `yaml:"normal,omitempty"`
	TemporaryValues Moment `yaml:"temporary,omitempty"`
	Metrics []StructuredData `yaml:"metrics,omitempty"`
	Moment Moment `yaml:"moment,omitempty"`
}

type ZwaveParams struct {
	Name string `yaml:"name,omitempty"`
	Ip string `yaml:"ip,omitempty"`
	Id string `yaml:"id,omitempty"`
	CommandClass string `yaml:"commandClass,omitempty"`
	Instance string `yaml:"instance,omitempty"`
	Value float64 `yaml:"value,omitempty"`
}

type URLS struct {
	Url string `yaml:"url"`
}

type StructuredData struct {
	Metric string `yaml:"metric,omitempty"`
	Labels map[string]string `yaml:"Labels,omitempty"`
	Timestamp string `yaml:"Timestamp,omitempty"`
	Timestamp2 string `yaml:"timestamp2,omitempty"`
	Value string `yaml:"temporary,omitempty"`
}

type GlobalSettingsType struct {
	Levels []Level `yaml:"levels,omitempty"`
	ActualHeater ZwaveParams `yaml:"heaterMetric,omitempty"`
	ActualTemperature ZwaveParams `yaml:"temperatureMetric,omitempty"`
	ApplicationRunningPath string `yaml:"applicationRunningPath,omitempty"`
	AuthorizedLevels []string `yaml:"authorizedLevels,omitempty"`
	LastUpdate time.Time `yaml:"lastUpdate,omitempty"`
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
	Moment time.Time `yaml:"moment,omitempty"`
	Time int `yaml:"time,omitempty"`
	Weekday time.Weekday `yaml:"weekday,omitempty"`
	Date int `yaml:"date,omitempty"`
	Level string `yaml:"level,omitempty"`
}

type Status struct {
	Heater_Level float64 `yaml:"level,omitempty"`
	Temperature_Requested float64 `yaml:"temperature_requested,omitempty"`
	Temperature_Actual float64 `yaml:"temperature_actual,omitempty"`
	Until time.Time `yaml:"until,omitempty"`
	TemporaryLevel string `yaml:"temporaryLevel,omitempty"`
	IsTemporary bool `yaml:"isTemporary,omitempty"`
	IpPort string `yaml:"ipPort,omitempty"`
	UpdateTime time.Time `yaml:"updateTime,omitempty"`
}