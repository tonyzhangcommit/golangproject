package config


type App struct {
	Env string `mapstructure:"env"`
	Port string `mapstructure:"port"`
	AppName string `mapstructure:"app_name"`
	AppUrl string `mapstructure:"app_url"`
}