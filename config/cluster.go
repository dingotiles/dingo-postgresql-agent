package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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

func FetchClusterSpec() (cluster *ClusterSpecification, err error) {
	apiSpec := APISpec()
	apiClusterSpec := fmt.Sprintf("%s/api", apiSpec.APIURI)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(apiClusterSpec)
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
