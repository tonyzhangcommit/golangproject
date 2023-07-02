package config

type Configuration struct {
	App App `mapstructure:"app"`
	Log Log `mapstructure:"log"`
	Database Database `mapstructure:"database"`
}
