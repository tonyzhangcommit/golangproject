package bootstrap

import (
	"orderingsystem/global"
	"orderingsystem/utils"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	level  zapcore.Level
	option []zap.Option
	corArr []zapcore.Core
)

func InitializeLog() *zap.Logger {
	utils.CreateDir(global.App.Config.Log.RootDir, os.ModePerm)

	if global.App.Config.Log.ShowLine {
		option = append(option, zap.AddCaller())
	}
	encoder := getEncoder()

	// 自定义日志级别
	highPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < zap.ErrorLevel && l >= zap.InfoLevel
	})

	DBPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.InfoLevel
	})

	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(getInfoLogWriter(), zapcore.AddSync(os.Stdout)), lowPriority)
	errorCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(getErrorLogWriter(), zapcore.AddSync(os.Stdout)), highPriority)
	DBCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(getDatabaseWriter(), zapcore.AddSync(os.Stdout)), DBPriority)

	corArr = append(corArr, infoCore)
	corArr = append(corArr, errorCore)
	corArr = append(corArr, DBCore)

	global.App.Log = zap.New(zapcore.NewTee(corArr...), option...)

	return zap.New(zapcore.NewTee(corArr...), option...)

}

func getEncoder() (encoder zapcore.Encoder) {
	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(global.App.Config.App.Env + "." + l.String())
	}

	// 设置编码器
	if global.App.Config.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return
}

func getInfoLogWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Log.InfoLog,
		MaxSize:    global.App.Config.Log.MaxSize,
		MaxBackups: global.App.Config.Log.MaxBackups,
		MaxAge:     global.App.Config.Log.MaxAge,
		Compress:   global.App.Config.Log.Compress,
	}
	return zapcore.AddSync(file)
}

func getErrorLogWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Log.ErrorLog,
		MaxSize:    global.App.Config.Log.MaxSize,
		MaxBackups: global.App.Config.Log.MaxBackups,
		MaxAge:     global.App.Config.Log.MaxAge,
		Compress:   global.App.Config.Log.Compress,
	}
	return zapcore.AddSync(file)
}

func getDatabaseWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Database.LogFilename,
		MaxSize:    global.App.Config.Log.MaxSize,
		MaxBackups: global.App.Config.Log.MaxBackups,
		MaxAge:     global.App.Config.Log.MaxAge,
		Compress:   global.App.Config.Log.Compress,
	}
	return zapcore.AddSync(file)
}
