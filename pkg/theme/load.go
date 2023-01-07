package theme

import (
	"github.com/spf13/viper"
)

func LoadTheme(v *viper.Viper, first bool) (*Theme, error) {
	theme := DefaultTheme()
	err := v.UnmarshalKey("theme", theme)
	if err != nil {
		panic("failed to unmarshal theme section from configuration file")
	}
	if !first || theme.File == "" {
		return theme, nil
	}

	v = viper.New()
	v.SetConfigFile(theme.File)
	err = v.ReadInConfig()
	if err != nil {
		return theme, err
	}
	return LoadTheme(v, false)
}
