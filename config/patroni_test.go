package config

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func patroniEnvForTesting() {
	os.Clearenv()
	os.Setenv("DINGO_IMAGE_VERSION", "0.0.0")
	os.Setenv("DINGO_CLUSTER", "unittest")
	os.Setenv("DINGO_ORG_TOKEN", "dummy")
	os.Setenv("DINGO_API_URI", "localhost:3000")
	os.Setenv("DINGO_PATRONI_DEFAULT_PATH", "../config/patroni-default-values.yml")
	os.Setenv("DOCKER_HOST_IP", "10.11.12.13")
	os.Setenv("DOCKER_HOST_PORT_5432", "5000")
	os.Setenv("DOCKER_HOST_PORT_8008", "8000")
}

func TestPatroni_CreateURIFile(t *testing.T) {
	patroniEnvForTesting()

	patroniSpec := &PatroniV12Specification{}
	patroniSpec.Postgresql.Authentication.Superuser.Username = "username"
	patroniSpec.Postgresql.Authentication.Superuser.Password = "password"
	uriFile := "/tmp/test_patroni/config/uri"
	err := patroniSpec.CreateURIFile(uriFile)
	if err != nil {
		t.Fatalf("Error creating /uri file: %s", err)
	}
	uri, err := ioutil.ReadFile(uriFile)
	if err != nil {
		t.Fatalf("Error reading %s file: %s", uriFile, err)
	}
	expectedURI := "postgres://username:password@10.11.12.13:5000/postgres"
	if string(uri) != expectedURI {
		t.Fatalf("URI %s did not match expected %s", uri, expectedURI)
	}
}

func TestPatroni_BasicPatroniSpec(t *testing.T) {
	patroniEnvForTesting()

	hostDiscoverySpec := &HostDiscoverySpecification{}
	clusterSpec := &ClusterSpecification{}
	overrideSpec := &OverrideSpecification{}

	clusterSpec.Etcd.URI = "https://localhost:4001"
	spec, err := BuildPatroniSpec(clusterSpec, hostDiscoverySpec, overrideSpec)
	if err != nil {
		t.Fatalf("Error creating patroni spec: %s", err)
	}

	if spec.Etcd.URL != clusterSpec.Etcd.URI {
		t.Fatalf("Etcd.URL not setup")
	}

	if spec.Bootstrap.Dcs.Postgresql.RecoveryConf.RestoreCommand != "/scripts/restore_command.sh \"%p\" \"%f\"" {
		t.Fatalf("recovery.conf's restore command should point to restore_command.sh")
	}

	if spec.Bootstrap.Dcs.Postgresql.RecoveryConf.TargetTimeline != "" {
		t.Fatalf("recovery.conf's recovery_target_timeline should be empty by default")
	}

	if strings.Contains(spec.String(), "target_timeline") {
		t.Fatalf("patroni.yml should not contain target_timeline if its not set")
	}
}

func TestPatroni_PatroniSpec_CustomRecoverySettings(t *testing.T) {
	patroniEnvForTesting()
	os.Setenv("PG_RECOVERY_TARGET_TIMELINE", "latest")

	hostDiscoverySpec := &HostDiscoverySpecification{}
	clusterSpec := &ClusterSpecification{}
	overrideSpec := OverrideSpec()
	spec, err := BuildPatroniSpec(clusterSpec, hostDiscoverySpec, overrideSpec)
	if err != nil {
		t.Fatalf("Error creating patroni spec: %s", err)
	}

	if spec.Bootstrap.Dcs.Postgresql.RecoveryConf.RestoreCommand != "/scripts/restore_command.sh \"%p\" \"%f\"" {
		t.Fatalf("recovery.conf's restore command should point to restore_command.sh")
	}

	if spec.Bootstrap.Dcs.Postgresql.RecoveryConf.TargetTimeline != "latest" {
		t.Fatalf("recovery.conf's recovery_target_timeline should be 'latest' from $PG_RECOVERY_TARGET_TIMELINE")
	}

	if !strings.Contains(spec.String(), "target_timeline") {
		t.Fatalf("patroni.yml should contain target_timeline if its not set")
	}
}
