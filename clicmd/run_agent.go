package clicmd

import (
	"fmt"
	"os"
	"path"
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
		fmt.Printf("Error trying to connect to API %s, retrying...\n", config.APISpec().APIURI)
		time.Sleep(time.Second)
		retryCount++
	}
	if err != nil {
		panic(err)
	}
	if clusterSpec == nil {
		fmt.Println("Cannot connect to API", config.APISpec().APIURI)
		os.Exit(1)
	}
	err = createPatroniPostgresConfigFiles(clusterSpec, "/", "postgres")
	if err != nil {
		panic(err)
	}

	startLongRunningAgent()
}

func createPatroniPostgresConfigFiles(clusterSpec *config.ClusterSpecification, rootPath string, postgresUser string) (err error) {
	fmt.Println(*clusterSpec)

	waleEnvDir := path.Join(rootPath, "/etc/wal-e.d/env")
	err = os.RemoveAll(waleEnvDir)
	if err != nil {
		return
	}
	environ := config.NewEnvironFromStrings(clusterSpec.WaleEnv)
	environ.AddEnv(fmt.Sprintf("REPLICATION_USER=%s", clusterSpec.Postgresql.Appuser.Username))
	environ.AddEnv("PG_DATA_DIR=/data/postgres0")

	err = environ.CreateEnvDirFiles(waleEnvDir)
	if err != nil {
		return
	}
	err = environ.CreateEnvScript(path.Join(rootPath, "/etc/patroni.d/.envrc"), postgresUser)
	if err != nil {
		return
	}

	patroniSpec, err := config.BuildPatroniSpec(clusterSpec, config.HostDiscoverySpec())
	if err != nil {
		return
	}
	err = patroniSpec.CreateConfigFile(path.Join(rootPath, "/config/patroni.yml"))
	return
}

func startLongRunningAgent() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true, // Output human readable JSON
	}))
	m.Get("/health", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"health": "ok"})
	})
	m.Run()
}
