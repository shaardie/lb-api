package server

import (
	"github.com/labstack/echo/v4"
)

type Server struct{}

func (*Server) GetHealth(ctx echo.Context) error {
	return nil
}

func (*Server) GetLoadbalancer(ctx echo.Context) error {
	return nil
}
