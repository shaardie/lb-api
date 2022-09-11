package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shaardie/lb-api/pkg/config"
	"github.com/shaardie/lb-api/pkg/db"
	"github.com/shaardie/lb-api/pkg/generate"
)

type Server struct {
	DB db.DB
}

func (*Server) GetHealth(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}

func (s *Server) GetLoadbalancers(ctx echo.Context) error {
	ct, err := s.DB.GetLoadBalancers(ctx.Request().Context())
	if err != nil {
		return s.sendError(ctx, http.StatusInternalServerError, "server error")
	}
	lbs := make([]generate.Loadbalancer, 0, len(ct))
	for _, lb := range ct {
		lbs = append(lbs, lb)
	}
	return ctx.JSON(http.StatusOK, lbs)
}

func (s *Server) GetLoadbalancer(ctx echo.Context, name string) error {
	lb, err := s.DB.GetLoadbalancer(ctx.Request().Context(), name)
	if err != nil {
		if err == db.ErrNotFound {
			return s.sendError(ctx, http.StatusNotFound, "not found")
		}
		return s.sendError(ctx, http.StatusInternalServerError, "server error")
	}
	return ctx.JSON(http.StatusOK, lb)
}

func (s *Server) CreateLoadBalancer(ctx echo.Context, name string) error {
	lb := generate.Loadbalancer{}
	err := ctx.Bind(&lb)
	if err != nil {
		return s.sendError(ctx, http.StatusBadRequest, "failed to parse input")
	}
	lb.Name = &name
	lb.Status = &generate.Status{
		Hostname: config.Cfg.Hostname,
		Ip:       config.Cfg.IP,
	}

	err = s.DB.CreateLoadbalancer(ctx.Request().Context(), name, lb)
	if err != nil {
		return s.sendError(ctx, http.StatusInternalServerError, "server error")
	}

	return ctx.JSON(http.StatusCreated, lb)
}

func (*Server) sendError(ctx echo.Context, code int, message string) error {
	genErr := generate.Error{
		Code:    code,
		Message: message,
	}
	return ctx.JSON(code, genErr)
}
