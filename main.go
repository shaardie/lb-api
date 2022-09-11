package main

import (
	"embed"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shaardie/lb-api/pkg/config"
	"github.com/shaardie/lb-api/pkg/db"
	"github.com/shaardie/lb-api/pkg/generate"
	"github.com/shaardie/lb-api/pkg/server"
)

//go:embed dist/**
var ui embed.FS

var (
	cfgHostname   = "loadbalancer.example.de"
	cfgDBFilename = "db.json"
)

func main() {
	config.Cfg = config.Config{
		DBFilename: cfgDBFilename,
		Hostname:   &cfgHostname,
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	generate.RegisterHandlers(e, &server.Server{DB: db.New(config.Cfg.DBFilename)})

	e.StaticFS("/ui", echo.MustSubFS(ui, "dist"))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
