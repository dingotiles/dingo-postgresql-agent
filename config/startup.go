package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type APISpecification struct {
	URI                string `required:"true"`
	PatroniDefaultPath string `default:"/patroni/patroni-default-values.yml" envconfig:"patroni_default_path"`
}

var apiSpec *APISpecification

func APISpec() *APISpecification {
	if apiSpec == nil {
		apiSpec = &APISpecification{}
		err := envconfig.Process("dingo_startup", apiSpec)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return apiSpec
}
