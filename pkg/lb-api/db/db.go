package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/shaardie/lb-api/pkg/generate"
	"github.com/shaardie/lb-api/pkg/lb-api/config"
	"github.com/shaardie/lb-api/pkg/lb-api/configurator"
)

var (
	ErrNotFound = errors.New("not found")
)

type DB interface {
	GetLoadBalancers(ctx context.Context) (Content, error)
	GetLoadbalancer(ctx context.Context, name string) (generate.Loadbalancer, error)
	CreateLoadbalancer(ctx context.Context, name string, cfg generate.Loadbalancer) error
	DeleteLoadbalancer(ctx context.Context, name string) error
}

type DBImpl struct {
	cfg      *config.Config
	cfgrator *configurator.Configurator
	m        *sync.Mutex
}

func New(cfg *config.Config, cfgrator *configurator.Configurator) DB {
	return &DBImpl{
		cfg:      cfg,
		cfgrator: cfgrator,
		m:        &sync.Mutex{},
	}
}

type Content map[string]generate.Loadbalancer

func (dbImpl *DBImpl) GetLoadBalancers(ctx context.Context) (Content, error) {
	dbImpl.m.Lock()
	defer dbImpl.m.Unlock()

	ct, err := dbImpl.read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read from database, %w", err)
	}
	return ct, nil
}

func (dbImpl *DBImpl) GetLoadbalancer(ctx context.Context, name string) (generate.Loadbalancer, error) {
	dbImpl.m.Lock()
	defer dbImpl.m.Unlock()

	ct, err := dbImpl.read(ctx)
	if err != nil {
		return generate.Loadbalancer{}, fmt.Errorf("failed to read from database, %w", err)
	}
	lb, ok := ct[name]
	if !ok {
		return generate.Loadbalancer{}, ErrNotFound
	}
	return lb, nil
}

func (dbImpl *DBImpl) CreateLoadbalancer(ctx context.Context, name string, cfg generate.Loadbalancer) error {
	dbImpl.m.Lock()
	defer dbImpl.m.Unlock()

	ct, err := dbImpl.read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read from database, %w", err)
	}

	ct[name] = cfg

	err = dbImpl.cfgrator.UpdateConfiguration(ctx, ct)
	if err != nil {
		return fmt.Errorf("configurator failed to update, %w", err)
	}

	err = dbImpl.write(ctx, ct)
	if err != nil {
		return fmt.Errorf("failed to update database, %w", err)
	}
	return nil
}

func (dbImpl *DBImpl) DeleteLoadbalancer(ctx context.Context, name generate.Name) error {
	dbImpl.m.Lock()
	defer dbImpl.m.Unlock()

	ct, err := dbImpl.read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read from database, %w", err)
	}

	delete(ct, name)

	err = dbImpl.cfgrator.UpdateConfiguration(ctx, ct)
	if err != nil {
		return fmt.Errorf("configurator failed to update, %w", err)
	}

	err = dbImpl.write(ctx, ct)
	if err != nil {
		return fmt.Errorf("failed to update database, %w", err)
	}
	return nil
}

func (dbImpl *DBImpl) read(ctx context.Context) (Content, error) {
	ct := make(Content)
	f, err := os.Open(dbImpl.cfg.DBFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return ct, nil
		}
		return ct, fmt.Errorf("unable to open file %v, %w", dbImpl.cfg.DBFilename, err)
	}

	jd := json.NewDecoder(f)
	err = jd.Decode(&ct)
	if err != nil {
		return ct, fmt.Errorf("failed to decode content from %v, %w", dbImpl.cfg.DBFilename, err)
	}

	return ct, nil
}

func (dbImpl *DBImpl) write(ctx context.Context, ct Content) error {
	f, err := os.Create(dbImpl.cfg.DBFilename)
	if err != nil {
		return fmt.Errorf("unable to override file %v, %w", dbImpl.cfg.DBFilename, err)
	}

	je := json.NewEncoder(f)
	je.SetIndent("", "  ")
	err = je.Encode(ct)
	if err != nil {
		return fmt.Errorf("failed to encode content to %v, %w", dbImpl.cfg.DBFilename, err)
	}

	return nil
}
