package config

type Log struct {
	RootDir    string `mapstructure:"root_dir"`
	InfoLog    string `mapstructure:"infolog"`
	ErrorLog   string `mapstructure:"errorlog"`
	ErrorInfoLog   string `mapstructure:"errorinfolog"`
	Format     string `mapstructure:"format"`
	ShowLine   bool   `mapstructure:"show_line"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}
