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
	Etcd struct {
		URI      string `json:"uri"`
		Host     string `json:"host"`
		Port     int16  `json:"port"`
		Protocol string `json:"protocol"`
		Username string `json:"username"`
		Password string `json:"password"`
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
	WaleEnv       []string `json:"wale_env"`
	RsyncArchives struct {
		Hostname             string `json:"hostname"`
		Username             string `json:"username"`
		SSHPort              string `json:"ssh_port"`
		DestinationDirectory string `json:"dest_dir"`
		PrivateKey           string `json:"private_key"`
	} `json:"rsync_archives"`
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
	return len(cluster.WaleEnv) > 0
}

// UsingRsync summarizes whether central API wants patroni to be configured for rysnc ship backups/WAL to remote host
func (cluster *ClusterSpecification) UsingRsync() bool {
	return cluster.RsyncArchives.Hostname != ""
}
