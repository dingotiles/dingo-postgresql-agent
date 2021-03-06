package clicmd

import (
	"fmt"
	"os"

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
		missingRequiredEnvs = []string{}

		clusterSpec := config.ClusterSpecification{}
		clusterSpec.Cluster.Name = req.NodeName
		clusterSpec.Cluster.Scope = req.ClusterName
		clusterSpec.Cluster.Namespace = requiredEnv("TEST_API_NAMESPACE")

		if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
			clusterSpec.Archives.Method = "s3"
			clusterSpec.Archives.S3.AWSAccessKeyID = requiredEnv("AWS_ACCESS_KEY_ID")
			clusterSpec.Archives.S3.AWSSecretAccessID = requiredEnv("AWS_SECRET_ACCESS_KEY")
			clusterSpec.Archives.S3.S3Bucket = requiredEnv("WAL_S3_BUCKET")
			clusterSpec.Archives.S3.S3Endpoint = requiredEnv("WALE_S3_ENDPOINT")
		} else if os.Getenv("SSH_HOST") != "" {
			clusterSpec.Archives.Method = "ssh"
			clusterSpec.Archives.SSH.Host = requiredEnv("SSH_HOST")
			clusterSpec.Archives.SSH.Port = requiredEnv("SSH_PORT")
			clusterSpec.Archives.SSH.User = requiredEnv("SSH_USER")
			clusterSpec.Archives.SSH.BasePath = requiredEnv("SSH_BASE_PATH")
			clusterSpec.Archives.SSH.PrivateKey = requiredEnv("SSH_PRIVATE_KEY")
		} else if os.Getenv("LOCAL_BACKUP_VOLUME") != "" {
			clusterSpec.Archives.Method = "local"
			clusterSpec.Archives.Local.LocalBackupVolume = requiredEnv("LOCAL_BACKUP_VOLUME")
		} else {
			missingRequiredEnvs = append(missingRequiredEnvs, "AWS_ACCESS_KEY_ID or SSH_HOST")
		}

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

func requiredEnv(envKey string) string {
	if os.Getenv(envKey) == "" {
		missingRequiredEnvs = append(missingRequiredEnvs, envKey)
	}
	return os.Getenv(envKey)
}
