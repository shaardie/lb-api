package loadbalancer

import (
	"context"
	"errors"
	"fmt"

	"github.com/shaardie/lb-api/pkg/generate"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

type LoadBalancer struct {
	client generate.ClientWithResponsesInterface
}

func (lb *LoadBalancer) apiStatusToServiceStatus(apiStatus *generate.Status) (status *v1.LoadBalancerStatus) {
	ingress := v1.LoadBalancerIngress{}
	if apiStatus != nil {
		if apiStatus.Ip != nil {
			ingress.IP = *apiStatus.Ip
		}
		if apiStatus.Hostname != nil {
			ingress.Hostname = *apiStatus.Hostname
		}
	}
	return &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{ingress}}
}

func (lb *LoadBalancer) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	return fmt.Sprintf("%s-%s-%s", clusterName, service.Namespace, service.Name)
}

func (lb *LoadBalancer) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	name := lb.GetLoadBalancerName(ctx, clusterName, service)
	resp, err := lb.client.GetLoadbalancerWithResponse(ctx, name)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get loadbalancer, %w", err)
	}
	if resp.JSON200 == nil {
		return nil, false, nil
	}

	return lb.apiStatusToServiceStatus(resp.JSON200.Status), true, nil
}

func (lb *LoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	name := lb.GetLoadBalancerName(ctx, clusterName, service)
	glb := generate.Loadbalancer{
		Config: generate.Config{
			Frontends: []generate.Frontend{},
		},
	}
	for _, port := range service.Spec.Ports {

		server := []string{}
		for _, node := range nodes {
			// Get Node Address
			addrs := node.Status.Addresses
			if len(addrs) == 0 {
				continue
			}
			server = append(server, node.Status.Addresses[0].Address)
		}
		glb.Config.Frontends = append(
			glb.Config.Frontends,
			generate.Frontend{
				Port: int(port.Port),
				Backend: generate.Backend{
					Server: &server,
				},
			},
		)
	}

	resp, err := lb.client.CreateLoadBalancerWithResponse(ctx, name, glb)
	if err != nil {
		return nil, fmt.Errorf("failed to call api and create load balancer, %w", err)
	}

	if resp.JSON201 == nil {
		return nil, errors.New("loadbalancer not created via API")
	}

	return lb.apiStatusToServiceStatus(resp.JSON201.Status), nil

}

func (lb *LoadBalancer) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	name := lb.GetLoadBalancerName(ctx, clusterName, service)
	resp, err := lb.client.DeleteLoadBalancerWithResponse(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to delete loadbalancer, %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("failed to delete loadbalancer with status %v: %v", resp.StatusCode(), resp.Status())
	}
	return nil
}
func (lb *LoadBalancer) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	return cloudprovider.NotImplemented
}
