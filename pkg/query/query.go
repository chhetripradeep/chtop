package query

type ClickHouseQueries struct {
	File    string  `json:"file" toml:"file" yaml:"file"`
	Queries []Query `json:"queries" toml:"queries" yaml:"queries"`
}

func DefaultClickHouseQueries() ClickHouseQueries {
	return ClickHouseQueries{}
}

type Query struct {
	Name       string
	Sql        string
	Latest     float64
	Datapoints []float64
}

func NewQuery(name, sql string) Query {
	return Query{
		Name: name,
		Sql:  sql,
	}
}

func (q *Query) Update(value float64) {
	q.Latest = value
	q.Datapoints = append(q.Datapoints, value)
}
