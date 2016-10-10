package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v1"
)

// PatroniV11Specification is constructed based on ClusterSpecification provided by the API.
// It is converted to a patroni.yml and used by Patroni to configure & run PostgreSQL.
// The scheme is for Patroni v1.1
type PatroniV11Specification struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Scope     string `yaml:"scope"`
	Restapi   struct {
		ConnectAddress string `yaml:"connect_address"`
		Listen         string `yaml:"listen"`
	} `yaml:"restapi"`
	Etcd struct {
		Host string `yaml:"host"`
	} `yaml:"etcd"`
	Bootstrap struct {
		Dcs struct {
			LoopWait             int `yaml:"loop_wait"`
			MaximumLagOnFailover int `yaml:"maximum_lag_on_failover"`
			Postgresql           struct {
				Parameters  interface{} `yaml:"parameters"`
				UsePgRewind bool        `yaml:"use_pg_rewind"`
			} `yaml:"postgresql"`
			RetryTimeout int `yaml:"retry_timeout"`
			TTL          int `yaml:"ttl"`
		} `yaml:"dcs"`
		Initdb []interface{} `yaml:"initdb"`
		PgHba  []string      `yaml:"pg_hba"`
		Users  struct {
			Admin struct {
				Options  []string `yaml:"options"`
				Password string   `yaml:"password"`
			} `yaml:"admin"`
		} `yaml:"users"`
	} `yaml:"bootstrap"`
	Postgresql struct {
		Authentication struct {
			Replication struct {
				Password string `yaml:"password"`
				Username string `yaml:"username"`
			} `yaml:"replication"`
			Superuser struct {
				Password string `yaml:"password"`
				Username string `yaml:"username"`
			} `yaml:"superuser"`
		} `yaml:"authentication"`
		ConnectAddress string `yaml:"connect_address"`
		DataDir        string `yaml:"data_dir"`
		Listen         string `yaml:"listen"`
		Parameters     struct {
			UnixSocketDirectories string `yaml:"unix_socket_directories"`
		} `yaml:"parameters"`
		Pgpass string `yaml:"pgpass"`
	} `yaml:"postgresql"`
	Tags struct {
		Clonefrom     bool `yaml:"clonefrom"`
		Nofailover    bool `yaml:"nofailover"`
		Noloadbalance bool `yaml:"noloadbalance"`
	} `yaml:"tags"`
}

var defaultPatroniSpec *PatroniV11Specification

func BuildPatroniSpec(clusterSpec *ClusterSpecification) (patroniSpec *PatroniV11Specification, err error) {
	patroniSpec, err = DefaultPatroniSpec()
	if err != nil {
		return
	}
	patroniSpec.MergeClusterSpec(clusterSpec)
	return
}

func DefaultPatroniSpec() (*PatroniV11Specification, error) {
	if defaultPatroniSpec == nil {
		filename, err := filepath.Abs(APISpec().PatroniDefaultPath)
		if err != nil {
			return nil, err
		}
		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		defaultPatroniSpec = &PatroniV11Specification{}
		err = yaml.Unmarshal(yamlFile, defaultPatroniSpec)
		if err != nil {
			return nil, err
		}
	}
	return defaultPatroniSpec, nil
}

func (patroniSpec *PatroniV11Specification) MergeClusterSpec(clusterSpec *ClusterSpecification) {
	patroniSpec.Etcd.Host = clusterSpec.Etcd.URI
	patroniSpec.Scope = clusterSpec.Cluster.Scope
	patroniSpec.Name = clusterSpec.Cluster.Name
}

func (patroniSpec *PatroniV11Specification) String() string {
	bytes, err := yaml.Marshal(patroniSpec)
	if err != nil {
		panic(err)
	}
	return string(bytes[:])
}
