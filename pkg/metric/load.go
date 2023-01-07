package metric

import (
	"github.com/spf13/viper"
)

func LoadMetrics(v *viper.Viper, first bool) (ClickHouseMetrics, error) {
	metrics := DefaultClickHouseMetrics()
	err := v.UnmarshalKey("clickhousemetrics", &metrics)
	if err != nil {
		panic("failed to unmarshal metrics section from configuration file")
	}
	if !first || metrics.File == "" {
		return metrics, nil
	}

	v = viper.New()
	v.SetConfigFile(metrics.File)
	err = v.ReadInConfig()
	if err != nil {
		return metrics, err
	}
	return LoadMetrics(v, false)
}
