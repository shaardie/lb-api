package cloudprovider

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	cloudprovider "k8s.io/cloud-provider"

	"github.com/shaardie/lb-api/pkg/cloud-provider-manager/loadbalancer"
	"github.com/shaardie/lb-api/pkg/generate"
)

const (
	ProviderName = "lb-api"
)

type providerConfig struct {
	LoadBalancer struct {
		URL string `yaml:"url"`
	} `yaml:"loadbalancer"`
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {

		if config == nil {
			return nil, errors.New("no cloud config file")
		}
		dc := yaml.NewDecoder(config)

		cfg := &providerConfig{}
		err := dc.Decode(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode cloud config, %w", err)
		}

		cli, err := generate.NewClientWithResponses(cfg.LoadBalancer.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to create lb-api client, %w", err)
		}
		return LBApiCloudProvider{
			lb: &loadbalancer.LoadBalancer{Client: cli},
		}, nil
	})
}

type LBApiCloudProvider struct {
	lb *loadbalancer.LoadBalancer
}

func (LBApiCloudProvider) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}

func (cp LBApiCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return cp.lb, true
}
func (LBApiCloudProvider) Instances() (cloudprovider.Instances, bool)     { return nil, false }
func (LBApiCloudProvider) InstancesV2() (cloudprovider.InstancesV2, bool) { return nil, false }
func (LBApiCloudProvider) Zones() (cloudprovider.Zones, bool)             { return nil, false }
func (LBApiCloudProvider) Clusters() (cloudprovider.Clusters, bool)       { return nil, false }
func (LBApiCloudProvider) Routes() (cloudprovider.Routes, bool)           { return nil, false }
func (LBApiCloudProvider) ProviderName() string                           { return ProviderName }
func (LBApiCloudProvider) HasClusterID() bool                             { return true }
