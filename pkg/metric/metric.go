package metric

type ClickHouseMetrics struct {
	File    string   `json:"file" toml:"file" yaml:"file"`
	Metrics []Metric `json:"metrics" toml:"metrics" yaml:"metrics"`
}

func DefaultClickHouseMetrics() ClickHouseMetrics {
	return ClickHouseMetrics{
		Metrics: []Metric{},
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
