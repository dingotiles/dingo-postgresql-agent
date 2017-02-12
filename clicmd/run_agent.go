package clicmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codeskyblue/go-sh"
	"github.com/dingotiles/dingo-postgresql-agent/config"
	"github.com/go-martini/martini"
	"github.com/hashicorp/errwrap"
	"github.com/martini-contrib/render"
)

// RunAgent runs the agent which fetches credentials/configuration,
// configures & runs Patroni, which in turn configures & runs PostgreSQL
func RunAgent(c *cli.Context) {
	startupConfig := config.APISpec()
	fmt.Println(*startupConfig)
	fmt.Println(*config.HostDiscoverySpec())
	retryCount := 0
	var err error
	var clusterSpec *config.ClusterSpecification
	for retryCount < 3 {
		clusterSpec, err = config.FetchClusterSpec()
		if err == nil && clusterSpec != nil {
			break
		}
		fmt.Printf("Error trying to connect to API %s, retrying... (%s)\n", startupConfig.APIURI, err)
		time.Sleep(time.Second)
		retryCount++
	}
	if err != nil {
		panic(err)
	}
	if clusterSpec == nil {
		fmt.Println("Cannot connect to API", startupConfig.APIURI)
		os.Exit(1)
	}
	err = createPatroniPostgresConfigFiles(clusterSpec, "/", "postgres")
	if err != nil {
		panic(err)
	}

	err = startPatroniPostgres(startupConfig)
	if err != nil {
		panic(err)
	}

	startLongRunningAgent()
}

func createPatroniPostgresConfigFiles(clusterSpec *config.ClusterSpecification, rootPath string, postgresUser string) (err error) {
	fmt.Println(*clusterSpec)

	patroniSpec, err := config.BuildPatroniSpec(clusterSpec, config.HostDiscoverySpec())
	if err != nil {
		return errwrap.Wrapf("Cannot BuildPatroniSpec: {{err}}", err)
	}

	environ := config.NewPatroniEnvironFromClusterSpec(clusterSpec)

	fmt.Println("Configuring continuous archives via wal-e")
	waleEnvDir := path.Join(rootPath, "/etc/wal-e.d/env")
	err = environ.CreateEnvDirFiles(waleEnvDir)
	if err != nil {
		return errwrap.Wrapf("Cannot create /etc/wal-e.d/env files: {{err}}", err)
	}

	err = environ.CreateEnvScript(path.Join(rootPath, "/etc/patroni.d/.envrc"), postgresUser)
	if err != nil {
		return errwrap.Wrapf("Cannot create /etc/patroni.d/.envrc: {{err}}", err)
	}

	err = patroniSpec.CreateConfigFile(path.Join(rootPath, "/config/patroni.yml"))
	if err != nil {
		return errwrap.Wrapf("Cannot create patroni.yml config file: {{err}}", err)
	}
	err = patroniSpec.CreateURIFile(path.Join(rootPath, "/config/uri"))
	if err != nil {
		return errwrap.Wrapf("Cannot create files with easy access to URIs: {{err}}", err)
	}
	return
}

func startPatroniPostgres(apiSpec *config.APISpecification) (err error) {
	if apiSpec.PatroniPostgresStartCmd == "" {
		fmt.Println("Assuming patroni & backup processes already runnning. No $PATRONI_POSTGRES_START_COMMAND start command provided.")
		return
	}
	cmdParts := strings.Split(apiSpec.PatroniPostgresStartCmd, " ")
	err = sh.Command(cmdParts[0], cmdParts[1:]).Run()
	if err != nil {
		return errwrap.Wrapf("Failed to run start command: {{err}}", err)
	}
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
