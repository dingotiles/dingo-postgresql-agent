package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestPatroni_Spec(t *testing.T) {
	os.Clearenv()
	os.Setenv("DOCKER_HOST_IP", "10.11.12.13")
	os.Setenv("DOCKER_HOST_PORT_5432", "5000")

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
