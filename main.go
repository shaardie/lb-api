package main

import (
	"embed"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shaardie/lb-api/pkg/generate"
	"github.com/shaardie/lb-api/pkg/server"
)

//go:embed dist/**
var ui embed.FS

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	generate.RegisterHandlers(e, &server.Server{})

	e.StaticFS("/ui", echo.MustSubFS(ui, "dist"))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
