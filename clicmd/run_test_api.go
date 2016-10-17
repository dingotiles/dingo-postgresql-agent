package clicmd

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// RunTestAPI runs the a sample backend API for which the Agent can be
// developed against.
func RunTestAPI(c *cli.Context) {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true, // Output human readable JSON
	}))
	m.Get("/api", func(r render.Render) {
		staticResponse := map[string]interface{}{
			"cluster": map[string]interface{}{
				"name":  "patroni1",
				"scope": "test-cluster-scope",
			},
			"wale_env": getWaleEnvVars(),
			// Example:
			// 	AWS_ACCESS_KEY_ID=AWS_ACCESS_KEY_ID
			// 	AWS_SECRET_ACCESS_KEY=AWS_SECRET_ACCESS_KEY
			// 	WAL_S3_BUCKET=WAL_S3_BUCKET
			// 	WALE_S3_ENDPOINT=https+path://s3.amazonaws.com:443
			// 	WALE_S3_PREFIX=s3://${WAL_S3_BUCKET}/backups/test-cluster-scope/wal/
			"postgresql": map[string]interface{}{
				"admin": map[string]interface{}{
					"password": "admin-password",
				},
				"superuser": map[string]interface{}{
					"username": "superuser-username",
					"password": "superuser-password",
				},
				"appuser": map[string]interface{}{
					"username": "appuser-username",
					"password": "appuser-password",
				},
			},
			"etcd": map[string]interface{}{
				"uri": os.Getenv("ETCD_HOST_PORT"),
			},
		}
		r.JSON(200, staticResponse)
	})
	m.Run()
}

func getWaleEnvVars() []string {
	return getWaleEnvVarsFromList(os.Environ())
}

func getWaleEnvVarsFromList(environ []string) []string {
	waleEnvCount := 0
	walePrefixes := []string{"WAL", "AWS", "WABS", "GOOGLE", "SWIFT", "PATRONI", "ETCD", "CONSUL"}
	for _, envVar := range environ {
		for _, prefix := range walePrefixes {
			if strings.Index(envVar, prefix) == 0 && !strings.HasSuffix(envVar, "=") {
				waleEnvCount++
			}
		}
	}
	waleEnvVars := make([]string, waleEnvCount)
	waleEnvIndex := 0
	for _, envVar := range environ {
		for _, prefix := range walePrefixes {
			if strings.Index(envVar, prefix) == 0 && !strings.HasSuffix(envVar, "=") {
				waleEnvVars[waleEnvIndex] = envVar
				waleEnvIndex++
			}
		}
	}
	return waleEnvVars
}
