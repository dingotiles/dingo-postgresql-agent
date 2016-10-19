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

type ClusterSpecification struct {
	Cluster struct {
		Name  string `json:"name"`
		Scope string `json:"scope"`
	} `json:"cluster"`
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
	WaleEnv []string `json:"wale_env"`
}

// TODO: POST ClusterName & OrgAuthToken to API

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
