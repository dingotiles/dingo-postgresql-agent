package clicmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/dingotiles/dingo-postgresql-agent/config"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var missingRequiredEnvs = []string{}

// RunTestAPI runs the a sample backend API for which the Agent can be
// developed against.
func RunTestAPI(c *cli.Context) {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true, // Output human readable JSON
	}))

	m.Post("/wal-e/api", binding.Bind(config.ContainerStartupRequest{}), func(req config.ContainerStartupRequest, r render.Render) {
		fmt.Printf("Recv [wal-e]: container start request: %v\n", req)
		name := "patroni1"
		patroniScope := "test-cluster-scope"

		clusterSpec := config.ClusterSpecification{}
		clusterSpec.Cluster.Name = name
		clusterSpec.Cluster.Scope = patroniScope

		clusterSpec.Archives.Method = "wal-e"
		clusterSpec.Archives.WalE.AWSAccessKeyID = requiredEnv("AWS_ACCESS_KEY_ID")
		clusterSpec.Archives.WalE.AWSSecretAccessID = requiredEnv("AWS_SECRET_ACCESS_KEY")
		clusterSpec.Archives.WalE.S3Bucket = requiredEnv("WAL_S3_BUCKET")
		clusterSpec.Archives.WalE.S3Endpoint = requiredEnv("WALE_S3_ENDPOINT")

		clusterSpec.Etcd.URI = requiredEnv("ETCD_URI")

		clusterSpec.Postgresql.Admin.Password = "admin-password"
		clusterSpec.Postgresql.Superuser.Username = "superuser-username"
		clusterSpec.Postgresql.Superuser.Password = "superuser-password"
		clusterSpec.Postgresql.Appuser.Username = "appuser-username"
		clusterSpec.Postgresql.Appuser.Password = "appuser-password"

		if len(missingRequiredEnvs) != 0 {
			fmt.Println("Missing required env:", missingRequiredEnvs)
			r.JSON(500, map[string]interface{}{"missing-env": missingRequiredEnvs})
			return
		}

		r.JSON(200, clusterSpec)
	})

	m.Post("/rsync-backup/api", binding.Bind(config.ContainerStartupRequest{}), func(req config.ContainerStartupRequest, r render.Render) {
		fmt.Printf("Recv [rsync-backup]: container start request: %v\n", req)
		name := "patroni1"
		patroniScope := "rsync-backup-cluster-scope"

		clusterSpec := config.ClusterSpecification{}
		clusterSpec.Cluster.Name = name
		clusterSpec.Cluster.Scope = patroniScope

		clusterSpec.Archives.Method = "rsync"
		clusterSpec.Archives.Rsync.URI = requiredEnv("RSYNC_URI")

		clusterSpec.Etcd.URI = requiredEnv("ETCD_URI")

		clusterSpec.Postgresql.Admin.Password = "admin-password"
		clusterSpec.Postgresql.Superuser.Username = "superuser-username"
		clusterSpec.Postgresql.Superuser.Password = "superuser-password"
		clusterSpec.Postgresql.Appuser.Username = "appuser-username"
		clusterSpec.Postgresql.Appuser.Password = "appuser-password"

		if len(missingRequiredEnvs) != 0 {
			fmt.Println("Missing required env:", missingRequiredEnvs)
			r.JSON(500, map[string]interface{}{"missing-env": missingRequiredEnvs})
			return
		}

		r.JSON(200, clusterSpec)
	})

	m.Run()
}

var (
	waleEnvVarPrefixes  = []string{"WAL", "AWS", "WABS", "GOOGLE", "SWIFT", "PATRONI", "ETCD", "CONSUL"}
	rsyncEnvVarPrefixes = []string{"RSYNC", "PATRONI", "ETCD", "CONSUL"}
)

func filterWaleEnvVars() []string {
	return filterEnvVarsFromList(os.Environ(), waleEnvVarPrefixes)
}

func filterRsyncEnvVars() []string {
	return filterEnvVarsFromList(os.Environ(), rsyncEnvVarPrefixes)
}

func filterEnvVarsFromList(environ, envVarPrefixes []string) []string {
	envCount := 0
	for _, envVar := range environ {
		for _, prefix := range envVarPrefixes {
			if strings.Index(envVar, prefix) == 0 && !strings.HasSuffix(envVar, "=") {
				envCount++
			}
		}
	}
	envVars := make([]string, envCount)
	envIndex := 0
	for _, envVar := range environ {
		for _, prefix := range envVarPrefixes {
			if strings.Index(envVar, prefix) == 0 && !strings.HasSuffix(envVar, "=") {
				envVars[envIndex] = envVar
				envIndex++
			}
		}
	}
	return envVars
}

func requiredEnv(envKey string) string {
	if os.Getenv(envKey) == "" {
		missingRequiredEnvs = append(missingRequiredEnvs, envKey)
	}
	return os.Getenv(envKey)
}
