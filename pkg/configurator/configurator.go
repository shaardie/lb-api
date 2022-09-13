package configurator

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"text/template"

	"github.com/shaardie/lb-api/pkg/config"
	"github.com/shaardie/lb-api/pkg/generate"
)

//go:embed haproxy.tpl
var tplFs embed.FS

type Configurator struct {
	template *template.Template
	m        *sync.Mutex
	cfg      *config.Config
}

func New(cfg *config.Config) (*Configurator, error) {
	t, err := template.New("haproxy.tpl").ParseFS(tplFs, "haproxy.tpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template, %w", err)
	}
	return &Configurator{
		template: t,
		m:        &sync.Mutex{},
		cfg:      cfg,
	}, nil
}

func (cfgurator *Configurator) UpdateConfiguration(ctx context.Context, lbs map[string]generate.Loadbalancer) error {
	cfgurator.m.Lock()
	defer cfgurator.m.Unlock()

	buf := new(bytes.Buffer)
	err := cfgurator.template.Execute(buf, lbs)
	if err != nil {
		return fmt.Errorf("failed to execute template, %w", err)
	}
	wantedCfgContent := buf.Bytes()

	currentCfgContent, err := os.ReadFile(cfgurator.cfg.ConfiguratorFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read configurator filename %v, %w", cfgurator.cfg.ConfiguratorFilename, err)
		}
		currentCfgContent = []byte{}
	}

	if reflect.DeepEqual(currentCfgContent, wantedCfgContent) {
		return nil
	}

	err = os.WriteFile(cfgurator.cfg.ConfiguratorFilename, wantedCfgContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write configurator filename %v, %w", cfgurator.cfg.ConfiguratorFilename, err)
	}

	err = exec.CommandContext(
		ctx,
		cfgurator.cfg.ConfiguratorCommand[0],
		cfgurator.cfg.ConfiguratorCommand[1:]...).Run()
	if err != nil {
		return fmt.Errorf("failed to reload, %w", err)
	}

	return nil
}
