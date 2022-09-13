package main

import (
	"embed"
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shaardie/lb-api/pkg/config"
	"github.com/shaardie/lb-api/pkg/configurator"
	"github.com/shaardie/lb-api/pkg/db"
	"github.com/shaardie/lb-api/pkg/generate"
	"github.com/shaardie/lb-api/pkg/server"
)

//go:embed dist/**
var ui embed.FS

var configFilename = flag.String("config", "lba-api.yaml", "name of the configuration file of the lba-api")

func main() {
	cfg, err := config.New(*configFilename)
	if err != nil {
		panic(err)
	}

	cfgrator, err := configurator.New(cfg)
	if err != nil {
		panic(err)
	}

	db := db.New(cfg, cfgrator)

	s := server.New(cfg, db, cfgrator)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	generate.RegisterHandlers(e, s)

	e.StaticFS("/ui", echo.MustSubFS(ui, "dist"))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
