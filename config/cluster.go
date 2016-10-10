package config

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
