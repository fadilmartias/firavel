package logger

import (
	"io"
	"os"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Debug  func(args ...any)
	Debugf func(template string, args ...any)

	Info  func(args ...any)
	Infof func(template string, args ...any)

	Warn  func(args ...any)
	Warnf func(template string, args ...any)

	Error  func(args ...any)
	Errorf func(template string, args ...any)
)

func Init() {
	options := []rotatelogs.Option{
		rotatelogs.WithMaxAge(30 * 24 * time.Hour),  // keep for 30 days
		rotatelogs.WithRotationTime(24 * time.Hour), // rotate daily
	}

	// Only use symlink on non-Windows
	if runtime.GOOS != "windows" {
		options = append(options, rotatelogs.WithLinkName("./storage/logs/app.log"))
	}

	fileWriter, err := rotatelogs.New(
		"./storage/logs/app-%Y-%m-%d.log", // template filename
		options...,
	)
	if err != nil {
		panic(err)
	}

	// Encoder Config
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout("02-01-2006 15:04:05"),
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// Output destination
	consoleWriter := zapcore.AddSync(os.Stdout)
	fileZapWriter := zapcore.AddSync(io.MultiWriter(fileWriter))

	// Encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

	// Level selectors
	// highPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
	// 	return l >= zapcore.ErrorLevel
	// })
	// lowPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
	// 	return l < zapcore.ErrorLevel
	// })

	// Cores
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel)
	fileCore := zapcore.NewCore(jsonEncoder, fileZapWriter, zapcore.DebugLevel)

	// Combine cores
	core := zapcore.NewTee(
		consoleCore, // all levels to console
		fileCore,    // only error+ to file
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar := logger.Sugar()

	// Expose global functions
	Debug = sugar.Debug
	Debugf = sugar.Debugf

	Info = sugar.Info
	Infof = sugar.Infof

	Warn = sugar.Warn
	Warnf = sugar.Warnf

	Error = sugar.Error
	Errorf = sugar.Errorf
}
