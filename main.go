package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"demo_railway/config"
	"demo_railway/internals/factory"
	middleware "demo_railway/internals/middleware"
	"demo_railway/internals/pkg/database"
	httpserver "demo_railway/internals/server"
	"demo_railway/internals/utils/env"
)

func init() {
	selectedEnv := config.Env()
	env := env.NewEnv()
	env.Load(`.env`)
	log.Info().Msg("Choosen environment " + selectedEnv)
}

// @title demo_railway-Project
// @version 0.0.1
// @description This is a doc for demo_railway-Project

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
func main() {
	cfg := config.Get()

	port := cfg.App.Port

	logLevel, err := zerolog.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	database.Init("std")

	f := factory.NewFactory()
	// f.Db.AutoMigrate(&example_feat.UserModel{})

	e := echo.New()
	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPDirect()
	middleware.Init(e)
	httpserver.Init(e, f)

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}
