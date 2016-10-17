package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type APISpecification struct {
	ClusterName        string `required:"true" envconfig:"cluster"`
	OrgAuthToken       string `required:"true" envconfig:"org_token"`
	APIURI             string `required:"true" envconfig:"api_uri"`
	PatroniDefaultPath string `default:"/patroni/patroni-default-values.yml" envconfig:"patroni_default_path"`
}

var apiSpec *APISpecification

func APISpec() *APISpecification {
	if apiSpec == nil {
		apiSpec = &APISpecification{}
		err := envconfig.Process("dingo", apiSpec)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return apiSpec
}
