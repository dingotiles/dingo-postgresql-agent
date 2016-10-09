package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type APISpecification struct {
	URI string `required:"true"`
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
