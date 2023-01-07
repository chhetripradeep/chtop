package metric

type ClickHouseMetrics struct {
	File    string   `json:"file" toml:"file" yaml:"file"`
	Metrics []Metric `json:"metrics" toml:"metrics" yaml:"metrics"`
}

func DefaultClickHouseMetrics() ClickHouseMetrics {
	return ClickHouseMetrics{
		Metrics: []Metric{
			NewMetric("ClickHouseProfileEvents_Query"),
			NewMetric("ClickHouseProfileEvents_SelectQuery"),
			NewMetric("ClickHouseProfileEvents_InsertQuery"),
			NewMetric("ClickHouseMetrics_PartsActive"),
			NewMetric("ClickHouseMetrics_TCPConnection"),
		},
	}
}

type Metric struct {
	Name       string
	Alias      string
	Latest     float64
	Datapoints []float64
}

func NewMetric(name string) Metric {
	return Metric{
		Name: name,
	}
}

func (m *Metric) Update(value float64) {
	m.Latest = value
	m.Datapoints = append(m.Datapoints, value)
}
