package metric

type Metrics []Metric

type Metric struct {
	Name       string
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
