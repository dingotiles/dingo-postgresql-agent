package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type HostDiscoverySpecification struct {
	IP       string `required:"true" envconfig:"docker_host_ip"`
	Port5432 string `required:"true" envconfig:"docker_host_port_5432"`
	Port8008 string `required:"true" envconfig:"docker_host_port_8008"`
}

var hostDiscoverySpec *HostDiscoverySpecification

func HostDiscoverySpec() *HostDiscoverySpecification {
	if hostDiscoverySpec == nil {
		hostDiscoverySpec = &HostDiscoverySpecification{}
		err := envconfig.Process("docker_host", hostDiscoverySpec)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return hostDiscoverySpec
}
