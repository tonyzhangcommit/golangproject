package global

import (
	"orderingsystem/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	Config      config.Configuration
	Configviper *viper.Viper
	Log         *zap.Logger
	DB          *gorm.DB
}

var App = new(Application)
