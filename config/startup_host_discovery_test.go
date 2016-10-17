package config

import (
	"os"
	"testing"
)

func TestStartupHostDiscovery(t *testing.T) {
	t.Parallel()

	os.Clearenv()
	os.Setenv("DOCKER_HOST_IP", "10.11.12.13")
	os.Setenv("DOCKER_HOST_PORT_5432", "5000")
	os.Setenv("DOCKER_HOST_PORT_8008", "8000")

	spec := HostDiscoverySpec()
	if spec.IP != "10.11.12.13" {
		t.Fatalf("IP should be 10.11.12.13, got %s", spec.IP)
	}
	if spec.Port5432 != "5000" {
		t.Fatalf("Port5432 should be 5000, got %s", spec.Port5432)
	}
	if spec.Port8008 != "8000" {
		t.Fatalf("Port8008 should be 8000, got %s", spec.Port8008)
	}

	os.Clearenv()
	os.Setenv("DOCKER_HOST_IP", "10.11.12.13")
	hostDiscoverySpec = nil
	spec = HostDiscoverySpec()
	if spec.IP != "10.11.12.13" {
		t.Fatalf("IP should be 10.11.12.13, got %s", spec.IP)
	}
	if spec.Port5432 != "5432" {
		t.Fatalf("Port5432 should default to 5432, got %s", spec.Port5432)
	}
	if spec.Port8008 != "8008" {
		t.Fatalf("Port8008 should default to 8008, got %s", spec.Port8008)
	}
}
