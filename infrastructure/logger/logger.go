package logger

import (
	//"backend-golang/infrastructure/config"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	// Load the default production configuration
	zapConfig := zap.NewProductionConfig()

	// Customize the encoder configuration
	zapConfig.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	zapConfig.DisableCaller = true

	// Build the logger
	logger, err := zapConfig.Build()
	if err != nil {
		panic("Cannot initialize logger: " + err.Error())
	}

	logger.Info("Logger initialized", zap.String("log_level", logger.Level().String()))

	return logger
}

func CloseLogger(logger *zap.Logger) {
	err := logger.Sync()
	if err != nil {
		log.Errorf("logger closed with an error: %s", err)
	}
}
