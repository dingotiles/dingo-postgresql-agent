package clicmd

import (
	"os"
	"testing"

	"github.com/dingotiles/dingo-postgresql-agent/config"
)

func TestRunAgent_createPatroniPostgresConfigFilesForWale(t *testing.T) {
	os.Clearenv()
	os.Setenv("DINGO_IMAGE_VERSION", "0.0.1")
	os.Setenv("DINGO_CLUSTER", "test-cluster")
	os.Setenv("DINGO_ORG_TOKEN", "test-org-token")
	os.Setenv("DINGO_API_URI", "localhost:3000")
	os.Setenv("DINGO_PATRONI_DEFAULT_PATH", "../config/patroni-wale-default-values.yml")
	os.Setenv("DOCKER_HOST_IP", "10.11.12.13")
	os.Setenv("DOCKER_HOST_PORT_5432", "5000")
	os.Setenv("DOCKER_HOST_PORT_8008", "8000")

	clusterSpec := &config.ClusterSpecification{}
	clusterSpec.Postgresql.Appuser.Username = "appuser"
	clusterSpec.Archives.Method = "wal-e"

	err := createPatroniPostgresConfigFiles(clusterSpec, "/tmp/run_agent_test", "")
	if err != nil {
		t.Fatalf("createPatroniPostgresConfigFiles for wal-e should not error; returned %s", err)
	}
}

func TestRunAgent_createPatroniPostgresConfigFilesForRsync(t *testing.T) {
	os.Clearenv()
	os.Setenv("DINGO_IMAGE_VERSION", "0.0.1")
	os.Setenv("DINGO_CLUSTER", "test-cluster")
	os.Setenv("DINGO_ORG_TOKEN", "test-org-token")
	os.Setenv("DINGO_API_URI", "localhost:3000")
	os.Setenv("DINGO_PATRONI_DEFAULT_PATH", "../config/patroni-rsync-default-values.yml")
	os.Setenv("DOCKER_HOST_IP", "10.11.12.13")
	os.Setenv("DOCKER_HOST_PORT_5432", "5000")
	os.Setenv("DOCKER_HOST_PORT_8008", "8000")

	clusterSpec := &config.ClusterSpecification{}
	clusterSpec.Postgresql.Appuser.Username = "appuser"
	clusterSpec.Archives.Method = "rsync"
	clusterSpec.Archives.Rsync.URI = "rsync://localhost:873/volume"

	err := createPatroniPostgresConfigFiles(clusterSpec, "/tmp/run_agent_test", "")
	if err != nil {
		t.Fatalf("createPatroniPostgresConfigFiles for rsync should not error; returned %s", err)
	}
}

func TestRunAgent_cannotCreatePatroniConfigWithoutArchiveConfig(t *testing.T) {
	os.Clearenv()
	os.Setenv("DINGO_IMAGE_VERSION", "0.0.1")
	os.Setenv("DINGO_CLUSTER", "test-cluster")
	os.Setenv("DINGO_ORG_TOKEN", "test-org-token")
	os.Setenv("DINGO_API_URI", "localhost:3000")
	os.Setenv("DINGO_PATRONI_DEFAULT_PATH", "../config/patroni-rsync-default-values.yml")
	os.Setenv("DOCKER_HOST_IP", "10.11.12.13")
	os.Setenv("DOCKER_HOST_PORT_5432", "5000")
	os.Setenv("DOCKER_HOST_PORT_8008", "8000")

	clusterSpec := &config.ClusterSpecification{}
	err := createPatroniPostgresConfigFiles(clusterSpec, "/tmp/run_agent_test", "")
	if err == nil {
		t.Fatalf("Expected createPatroniPostgresConfigFiles to return error for missing wale_env and rsync_archives")
	}
}
