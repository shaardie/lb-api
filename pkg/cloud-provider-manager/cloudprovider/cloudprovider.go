package cloudprovider

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
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
		URL         string  `yaml:"url"`
		BearerToken string  `yaml:"bearer_token"`
		Certificate *string `yaml:"certificate"`
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

		bearerTokenProvider, bearerTokenProviderErr := securityprovider.NewSecurityProviderBearerToken(cfg.LoadBalancer.BearerToken)
		if bearerTokenProviderErr != nil {
			return nil, fmt.Errorf("failed to create bearer token provider, %w", err)
		}

		// Custom http.Client
		httpClient := &http.Client{}
		// Load custom certificate, if necessary
		if cfg.LoadBalancer.Certificate != nil {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(*cfg.LoadBalancer.Certificate))

			httpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			}
		}
		cli, err := generate.NewClientWithResponses(
			cfg.LoadBalancer.URL,
			generate.WithHTTPClient(httpClient),
			generate.WithRequestEditorFn(bearerTokenProvider.Intercept),
		)
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
