package config

type Database struct {
	Driver string `mapstructure:"driver"`
    Host string `mapstructure:"host"`
    Port int `mapstructure:"port"`
    Database string `mapstructure:"database"`
    UserName string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Charset string `mapstructure:"charset"`
    MaxIdleConns int `mapstructure:"max_idle_conns"`
    MaxOpenConns int `mapstructure:"max_open_conns"`
    LogMode string `mapstructure:"log_mode"`
    EnableFileLogWriter bool `mapstructure:"enable_file_log_writer"`
    LogFilename string `mapstructure:"log_filename"`
}