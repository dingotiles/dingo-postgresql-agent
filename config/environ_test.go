package config

import (
	"os"
	"testing"
)

func TestEnviron_Namespace(t *testing.T) {
	os.Clearenv()
	clusterSpec := &ClusterSpecification{}
	clusterSpec.Cluster.Name = "cluster1node1"
	clusterSpec.Cluster.Scope = "cluster1"
	clusterSpec.Cluster.Namespace = "/service_instances/xyz/service/"
	clusterSpec.Etcd.URI = "https://user:password@etcd.cluster:4001"
	environ := *NewPatroniEnvironFromClusterSpec(clusterSpec)
	if environ["ETCD_URI"] != "https://user:password@etcd.cluster:4001" {
		t.Fatalf("$ETCD_URI was '%s' should be 'https://user:password@etcd.cluster:4001'", environ["ETCD_URI"])
	}
	if environ["ETCD_CLUSTER_URI"] != "https://user:password@etcd.cluster:4001/v2/keys/service_instances/xyz/service/cluster1" {
		t.Fatalf("$ETCD_CLUSTER_URI was '%s' should be '%s'",
			environ["ETCD_CLUSTER_URI"],
			"https://user:password@etcd.cluster:4001/v2/keys/service_instances/xyz/service/cluster1")
	}
}
