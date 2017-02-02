package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ContainerStartupRequest is the expected inbound data when each Dingo Agent starts
type ContainerStartupRequest struct {
	ImageVersion string `json:"image_version"`
	ClusterName  string `json:"cluster"`
	OrgAuthToken string `json:"org_token"`
}

// ClusterSpecification describes the cluster configuration provided by central API
type ClusterSpecification struct {
	Cluster struct {
		Name  string `json:"name"`
		Scope string `json:"scope"`
	} `json:"cluster"`
	Archives struct {
		Method string `json:"method"`
		WalE   struct {
			AWSAccessKeyID    string `json:"aws_access_key_id,omitempty"`
			AWSSecretAccessID string `json:"aws_secret_access_id,omitempty"`
			S3Bucket          string `json:"s3_bucket,omitempty"`
			S3Endpoint        string `json:"s3_endpoint,omitempty"`
		} `json:"wale,omitempty"`
		Rsync struct {
			URI string `json:"uri,omitempty"`
		} `json:"rsync,omitempty"`
	} `json:"archives"`
	Etcd struct {
		URI string `json:"uri"`
	} `json:"etcd"`
	Postgresql struct {
		Admin struct {
			Password string `json:"password"`
		} `json:"admin"`
		Appuser struct {
			Password string `json:"password"`
			Username string `json:"username"`
		} `json:"appuser"`
		Superuser struct {
			Password string `json:"password"`
			Username string `json:"username"`
		} `json:"superuser"`
	} `json:"postgresql"`
}

// TODO: POST ClusterName & OrgAuthToken to API

// FetchClusterSpec retrieves the new/existing configuration for a cluster from central API
func FetchClusterSpec() (cluster *ClusterSpecification, err error) {
	apiSpec := APISpec()
	apiClusterSpec := fmt.Sprintf("%s/api", apiSpec.APIURI)
	fmt.Printf("Loading configuration from %s...", apiClusterSpec)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	startupReq := ContainerStartupRequest{
		ImageVersion: apiSpec.ImageVersion,
		ClusterName:  apiSpec.ClusterName,
		OrgAuthToken: apiSpec.OrgAuthToken,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(startupReq)

	resp, err := netClient.Post(apiClusterSpec, "application/json", b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &cluster)

	return
}

// UsingWale summarizes whether central API wants patroni to be configured for wal-e to ship backups/WAL
func (cluster *ClusterSpecification) UsingWale() bool {
	return cluster.Archives.Method == "wal-e"
}

// UsingRsync summarizes whether central API wants patroni to be configured for rysnc ship backups/WAL to remote host
func (cluster *ClusterSpecification) UsingRsync() bool {
	return cluster.Archives.Method == "rsync"
}

func (cluster *ClusterSpecification) waleS3Prefix() string {
	return fmt.Sprintf("s3://%s/backups/%s/wal/", cluster.Archives.WalE.S3Bucket, cluster.Cluster.Scope)
}
