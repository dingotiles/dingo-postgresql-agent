package clicmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dingotiles/dingo-postgresql-agent/config"
)

// RunAgent runs the agent which fetches credentials/configuration,
// configures & runs Patroni, which in turn configures & runs PostgreSQL
func RunAgent(c *cli.Context) {
	fmt.Printf("API config: %#v\n", *config.APISpec())
	config.FetchClusterSpec()
}
