package bootstrap

import (
	"fmt"
	"orderingsystem/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 这里初始化配置文件

func InitializeConfig() *viper.Viper {
	configPath := "config.yaml"
	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		configPath = configEnv
	}
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config file failed %s \n", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		if err := v.Unmarshal(&global.App.Config); err != nil {
			panic(fmt.Errorf("read config file failed %s \n", err))
		}
	})
	if err := v.Unmarshal(&global.App.Config); err != nil {
		panic(fmt.Errorf("read config file failed %s \n", err))
	}

	global.App.Configviper = v

	return v
}
