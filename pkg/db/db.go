package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/shaardie/lb-api/pkg/generate"
)

var (
	ErrNotFound = errors.New("not found")
)

type DB interface {
	GetLoadBalancers(ctx context.Context) (Content, error)
	GetLoadbalancer(ctx context.Context, name string) (generate.Loadbalancer, error)
	CreateLoadbalancer(ctx context.Context, name string, cfg generate.Loadbalancer) error
}

type DBImpl struct {
	filename string
	m        *sync.Mutex
}

func New(filename string) DB {
	return &DBImpl{
		filename: filename,
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

	err = dbImpl.write(ctx, ct)
	if err != nil {
		return fmt.Errorf("failed to update database, %w", err)
	}
	return nil
}

func (dbImpl *DBImpl) read(ctx context.Context) (Content, error) {
	ct := make(Content)
	f, err := os.Open(dbImpl.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return ct, nil
		}
		return ct, fmt.Errorf("unable to open file %v, %w", dbImpl.filename, err)
	}

	jd := json.NewDecoder(f)
	err = jd.Decode(&ct)
	if err != nil {
		return ct, fmt.Errorf("failed to decode content from %v, %w", dbImpl.filename, err)
	}

	return ct, nil
}

func (dbImpl *DBImpl) write(ctx context.Context, ct Content) error {
	f, err := os.Create(dbImpl.filename)
	if err != nil {
		return fmt.Errorf("unable to override file %v, %w", dbImpl.filename, err)
	}

	je := json.NewEncoder(f)
	je.SetIndent("", "  ")
	err = je.Encode(ct)
	if err != nil {
		return fmt.Errorf("failed to encode content to %v, %w", dbImpl.filename, err)
	}

	return nil
}
