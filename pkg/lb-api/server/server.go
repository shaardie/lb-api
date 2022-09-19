package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/shaardie/lb-api/pkg/generate"
	"github.com/shaardie/lb-api/pkg/lb-api/config"
	"github.com/shaardie/lb-api/pkg/lb-api/configurator"
	"github.com/shaardie/lb-api/pkg/lb-api/db"
)

type Server struct {
	DB  db.DB
	cfg *config.Config
}

func New(cfg *config.Config, db db.DB, configurator *configurator.Configurator) *Server {
	return &Server{
		DB:  db,
		cfg: cfg,
	}
}

func (*Server) GetHealth(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}

func (s *Server) GetLoadbalancers(ctx echo.Context) error {
	ct, err := s.DB.GetLoadBalancers(ctx.Request().Context())
	if err != nil {
		return err
	}
	lbs := make([]generate.Loadbalancer, 0, len(ct))
	for _, lb := range ct {
		lbs = append(lbs, lb)
	}
	return ctx.JSON(http.StatusOK, lbs)
}

func (s *Server) GetLoadbalancer(ctx echo.Context, name generate.Name) error {
	lb, err := s.DB.GetLoadbalancer(ctx.Request().Context(), name)
	if err != nil {
		if err == db.ErrNotFound {
			return s.sendError(ctx, http.StatusNotFound, "not found")
		}
		return err
	}
	return ctx.JSON(http.StatusOK, lb)
}

func (s *Server) CreateLoadBalancer(ctx echo.Context, name generate.Name) error {
	lb := generate.Loadbalancer{}
	err := ctx.Bind(&lb)
	if err != nil {
		return s.sendError(ctx, http.StatusBadRequest, "failed to parse input")
	}
	lb.Name = &name
	lb.Status = &generate.Status{
		Hostname: s.cfg.Hostname,
		Ip:       s.cfg.IP,
	}

	err = s.DB.CreateLoadbalancer(ctx.Request().Context(), name, lb)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, lb)
}

func (s *Server) DeleteLoadBalancer(ctx echo.Context, name generate.Name) error {
	err := s.DB.DeleteLoadbalancer(ctx.Request().Context(), name)
	if err != nil {
		return err
	}
	return nil
}

func (*Server) sendError(ctx echo.Context, code int, message string) error {
	genErr := generate.Error{
		Code:    code,
		Message: message,
	}
	return ctx.JSON(code, genErr)
}
