package main

import (
	"context"
	"crypto/subtle"
	"embed"
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shaardie/lb-api/pkg/generate"
	"github.com/shaardie/lb-api/pkg/lb-api/config"
	"github.com/shaardie/lb-api/pkg/lb-api/configurator"
	"github.com/shaardie/lb-api/pkg/lb-api/db"
	"github.com/shaardie/lb-api/pkg/lb-api/server"
)

//go:embed dist/**
var ui embed.FS

var configFilename = flag.String("config", "lb-api.yaml", "name of the configuration file of the lba-api")

func main() {
	flag.Parse()

	cfg, err := config.New(*configFilename)
	if err != nil {
		panic(err)
	}

	cfgrator, err := configurator.New(cfg)
	if err != nil {
		panic(err)
	}

	db := db.New(cfg, cfgrator)

	ctx := context.TODO()
	ct, err := db.GetLoadBalancers(ctx)
	if err != nil {
		panic(err)
	}
	err = cfgrator.UpdateConfiguration(ctx, ct)
	if err != nil {
		panic(err)
	}

	s := server.New(cfg, db, cfgrator)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register api with authentication
	api := e.Group("", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Validator: func(key string, c echo.Context) (bool, error) {
			return subtle.ConstantTimeCompare([]byte(key), []byte(cfg.BearerToken)) == 1, nil
		},
	}))
	generate.RegisterHandlers(api, s)

	// Serve UI
	e.StaticFS("/ui", echo.MustSubFS(ui, "dist"))

	// Start server with or without TLS
	if cfg.TLS == nil {
		e.Logger.Fatal(e.Start(cfg.AdminAddress))
	}
	e.Logger.Fatal(e.StartTLS(cfg.AdminAddress, cfg.TLS.CertificateFilename, cfg.TLS.KeyFilename))
}
