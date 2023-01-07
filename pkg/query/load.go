package query

import (
	"github.com/spf13/viper"
)

func LoadQueries(v *viper.Viper, first bool) (ClickHouseQueries, error) {
	queries := DefaultClickHouseQueries()
	err := v.UnmarshalKey("clickhousequeries", &queries)
	if err != nil {
		panic("failed to unmarshal Queries section from configuration file")
	}
	if !first || queries.File == "" {
		return queries, nil
	}

	v = viper.New()
	v.SetConfigFile(queries.File)
	err = v.ReadInConfig()
	if err != nil {
		return queries, err
	}
	return LoadQueries(v, false)
}
