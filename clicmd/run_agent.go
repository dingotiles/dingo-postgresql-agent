package clicmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dingotiles/dingo-postgresql-agent/config"
)

// RunAgent runs the agent which fetches credentials/configuration,
// configures & runs Patroni, which in turn configures & runs PostgreSQL
func RunAgent(c *cli.Context) {
	fmt.Println(*config.APISpec())
	clusterSpec, err := config.FetchClusterSpec()
	if err != nil {
		panic(err)
	}
	// fmt.Println(*clusterSpec)

	patroniSpec, err := config.BuildPatroniSpec(clusterSpec)
	if err != nil {
		panic(err)
	}
	fmt.Println(patroniSpec)
}
