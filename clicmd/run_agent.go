package clicmd

import (
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/dingotiles/dingo-postgresql-agent/config"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// RunAgent runs the agent which fetches credentials/configuration,
// configures & runs Patroni, which in turn configures & runs PostgreSQL
func RunAgent(c *cli.Context) {
	fmt.Println(*config.APISpec())
	fmt.Println(*config.HostDiscoverySpec())
	retryCount := 0
	var err error
	var clusterSpec *config.ClusterSpecification
	for retryCount < 3 {
		clusterSpec, err = config.FetchClusterSpec()
		if err == nil && clusterSpec != nil {
			break
		}
		fmt.Printf("Error trying to connect to API %s, retrying...\n", config.APISpec().URI)
		time.Sleep(time.Second)
		retryCount++
	}
	if err != nil {
		panic(err)
	}
	if clusterSpec == nil {
		fmt.Println("Cannot connect to API", config.APISpec().URI)
		os.Exit(1)
	}
	fmt.Println(*clusterSpec)

	patroniSpec, err := config.BuildPatroniSpec(clusterSpec, config.HostDiscoverySpec())
	if err != nil {
		panic(err)
	}
	fmt.Println(patroniSpec)

	envdir := config.NewEnvdirFromStrings(clusterSpec.WaleEnv)
	err = envdir.CreateFiles("/etc/wal-e.d/env")
	if err != nil {
		panic(err)
	}

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true, // Output human readable JSON
	}))
	m.Get("/health", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"health": "ok"})
	})
	m.Run()
}
