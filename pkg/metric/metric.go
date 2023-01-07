package metric

type ClickHouseMetrics struct {
	File    string   `json:"file" toml:"file" yaml:"file"`
	Metrics []Metric `json:"metrics" toml:"metrics" yaml:"metrics"`
}

func DefaultClickHouseMetrics() ClickHouseMetrics {
	return ClickHouseMetrics{
		Metrics: []Metric{
			NewMetric("ClickHouseProfileEvents_Query", "Total Queries"),
			NewMetric("ClickHouseProfileEvents_SelectQuery", "Total Select Queries"),
			NewMetric("ClickHouseProfileEvents_InsertQuery", "Total Insert Queries"),
			NewMetric("ClickHouseMetrics_PartsActive", "Number of Active Parts"),
			NewMetric("ClickHouseMetrics_TCPConnection", "Number of Open File Descriptors"),
		},
	}
}

type Metric struct {
	Name       string
	Alias      string
	Latest     float64
	Datapoints []float64
}

func NewMetric(name, alias string) Metric {
	return Metric{
		Name:  name,
		Alias: alias,
	}
}

func (m *Metric) Update(value float64) {
	m.Latest = value
	m.Datapoints = append(m.Datapoints, value)
}
