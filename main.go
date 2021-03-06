package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danymarita/static-files-server/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMapper = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func initLogger(cfg config.Provider) *zap.Logger {
	var level zapcore.Level
	if lvl, ok := levelMapper[cfg.GetString("logger.level")]; ok {
		level = lvl
	} else {
		level = zapcore.InfoLevel
	}

	loggerCfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339NanoTimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := loggerCfg.Build()
	return logger
}

func main() {
	e := echo.New()
	config := config.Config()
	logger := initLogger(config)
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		urlString := c.Request().URL.String()
		logger.Info("Intercept request",
			zap.String("method", c.Request().Method),
			zap.String("uri", urlString),
			zap.String("ip", c.Request().RemoteAddr),
		)
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	}))

	idleConnectionClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logger.Info("[API] Server is shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// We received an interrupt signal, shut down.
		if err := e.Shutdown(ctx); err != nil {
			logger.Error("[API] Fail to shutting down",
				zap.Error(err),
			)
		}
		close(idleConnectionClosed)
	}()

	fs := http.FileServer(http.Dir(config.GetString("folders.cms")))
	e.GET("/cms/*", echo.WrapHandler(http.StripPrefix("/cms/", fs)))

	srvAddr := fmt.Sprintf("%s:%d", config.GetString("app.host"), config.GetInt("app.port"))
	logger.Info(fmt.Sprintf("[API] HTTP serve at %s", srvAddr))
	if err := e.Start(srvAddr); err != nil {
		logger.Error("[API] Fail to start listen and server",
			zap.Error(err),
		)
	}

	<-idleConnectionClosed
	logger.Info("[API] Bye")
}
