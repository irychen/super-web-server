package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"super-web-server/internal/app"
	"super-web-server/internal/config"
	"super-web-server/internal/types"
	"super-web-server/pkg/logger"
	"syscall"
	"time"
)

const DefaultConfigPath = "./configs"

func main() {
	mode := ParseMode()
	config, err := config.LoadConfig(DefaultConfigPath, mode)

	if err != nil {
		panic(err)
	}

	if err := InitLogger(config); err != nil {
		panic(err)
	}

	app, err := app.NewApp(config)

	if err != nil {
		logger.Fatal(err.Error())
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()

	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
	}

	defer func() {
		logger.Sync()
	}()
}

func InitLogger(config *config.Config) error {
	var logConfig = config.Log

	level, err := logger.ParseStringLogLevel(logConfig.Level)
	if err != nil {
		return fmt.Errorf("parse log level failed: %w", err)
	}

	format := logger.FormatJSON
	if config.Mode == types.ServerModeDev {
		format = logger.FormatConsole
	}

	err = logger.InitLogger(logger.Config{
		Level:      level,
		Format:     format,
		Directory:  logConfig.Path,
		Filename:   logConfig.Filename,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		Compress:   logConfig.Compress,
		Stdout:     logConfig.Stdout,
	})

	if err != nil {
		return fmt.Errorf("init global logger failed: %w", err)
	}

	return nil
}

func ParseMode() types.ServerMode {
	var mode = "dev"
	var envMode = os.Getenv("APP_MODE")

	if envMode != "" {
		mode = envMode
	}

	var flagMode = flag.String("mode", mode, "set app mode, valid modes: dev, prod, test, local")

	flag.Parse()

	if *flagMode != "" {
		mode = *flagMode
	}

	serverMode, err := types.ParseServerMode(mode)
	if err != nil {
		panic(err)
	}

	return serverMode
}
