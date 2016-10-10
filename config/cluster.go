package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ClusterSpecification struct {
	Etcd struct {
		URI string `json:"uri"`
	} `json:"etcd"`
	Postgresql struct {
		Admin struct {
			Password string `json:"password"`
			Username string `json:"username"`
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
	WaleEnv  []string `json:"wale_env"`
	WaleMode string   `json:"wale_mode"`
}

func FetchClusterSpec() (cluster *ClusterSpecification, err error) {
	apiSpec := APISpec()
	apiClusterSpec := fmt.Sprintf("%s/api", apiSpec.URI)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(apiClusterSpec)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(body, &cluster)

	return
}
