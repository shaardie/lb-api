package cloudprovider

import (
	"io"

	cloudprovider "k8s.io/cloud-provider"

	"github.com/shaardie/lb-api/pkg/cloud-provider-manager/loadbalancer"
)

const (
	ProviderName = "lb-api"
)

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		return &LBApiCloudProvider{}, nil
	})
}

type LBApiCloudProvider struct {
	lb *loadbalancer.LoadBalancer
}

func (*LBApiCloudProvider) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}

func (cp *LBApiCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return cp.lb, true
}
func (*LBApiCloudProvider) Instances() (cloudprovider.Instances, bool)     { return nil, false }
func (*LBApiCloudProvider) InstancesV2() (cloudprovider.InstancesV2, bool) { return nil, false }
func (*LBApiCloudProvider) Zones() (cloudprovider.Zones, bool)             { return nil, false }
func (*LBApiCloudProvider) Clusters() (cloudprovider.Clusters, bool)       { return nil, false }
func (*LBApiCloudProvider) Routes() (cloudprovider.Routes, bool)           { return nil, false }
func (*LBApiCloudProvider) ProviderName() string                           { return ProviderName }
func (*LBApiCloudProvider) HasClusterID() bool                             { return false }
