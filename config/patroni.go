package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v1"
)

// PatroniV11Specification is constructed based on ClusterSpecification provided by the API.
// It is converted to a patroni.yml and used by Patroni to configure & run PostgreSQL.
// The scheme is for Patroni v1.1
type PatroniV11Specification struct {
	Scope     string `yaml:"scope"`
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Restapi   struct {
		Listen         string `yaml:"listen"`
		ConnectAddress string `yaml:"connect_address"`
	} `yaml:"restapi"`
	Etcd struct {
		Host string `yaml:"host"`
	} `yaml:"etcd"`
	Bootstrap struct {
		Dcs struct {
			TTL                  int `yaml:"ttl"`
			LoopWait             int `yaml:"loop_wait"`
			RetryTimeout         int `yaml:"retry_timeout"`
			MaximumLagOnFailover int `yaml:"maximum_lag_on_failover"`
			Postgresql           struct {
				UsePgRewind bool `yaml:"use_pg_rewind"`
				UseSlots    bool `yaml:"use_slots"`
				Parameters  struct {
					WalLevel            string `yaml:"wal_level"`
					HotStandby          string `yaml:"hot_standby"`
					WalKeepSegments     int    `yaml:"wal_keep_segments"`
					MaxWalSenders       int    `yaml:"max_wal_senders"`
					MaxReplicationSlots int    `yaml:"max_replication_slots"`
					WalLogHints         string `yaml:"wal_log_hints"`
					ArchiveMode         string `yaml:"archive_mode"`
					ArchiveTimeout      string `yaml:"archive_timeout"`
					ArchiveCommand      string `yaml:"archive_command"`
				} `yaml:"parameters"`
				RecoveryConf struct {
					RestoreCommand string `yaml:"restore_command"`
				} `yaml:"recovery_conf"`
			} `yaml:"postgresql"`
		} `yaml:"dcs"`
		Initdb []interface{} `yaml:"initdb"`
		PgHba  []string      `yaml:"pg_hba"`
		Users  struct {
			Postgres struct {
				Password string   `yaml:"password"`
				Options  []string `yaml:"options"`
			} `yaml:"postgres"`
		} `yaml:"users"`
	} `yaml:"bootstrap"`
	Postgresql struct {
		Listen         string `yaml:"listen"`
		ConnectAddress string `yaml:"connect_address"`
		DataDir        string `yaml:"data_dir"`
		Pgpass         string `yaml:"pgpass"`
		Authentication struct {
			Replication struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"replication"`
			Superuser struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"superuser"`
		} `yaml:"authentication"`
		Parameters struct {
			UnixSocketDirectories string `yaml:"unix_socket_directories"`
		} `yaml:"parameters"`
		Callbacks struct {
			OnStart      string `yaml:"on_start"`
			OnStop       string `yaml:"on_stop"`
			OnRestart    string `yaml:"on_restart"`
			OnRoleChange string `yaml:"on_role_change"`
		} `yaml:"callbacks"`
		CreateReplicaMethod []string `yaml:"create_replica_method"`
		WalE                struct {
			Command                       string `yaml:"command"`
			Envdir                        string `yaml:"envdir"`
			ThresholdMegabytes            int    `yaml:"threshold_megabytes"`
			ThresholdBackupSizePercentage int    `yaml:"threshold_backup_size_percentage"`
			Retries                       int    `yaml:"retries"`
			UseIam                        int    `yaml:"use_iam"`
			NoMaster                      int    `yaml:"no_master"`
		} `yaml:"wal_e"`
	} `yaml:"postgresql"`
	Tags struct {
		Nofailover    bool `yaml:"nofailover"`
		Noloadbalance bool `yaml:"noloadbalance"`
		Clonefrom     bool `yaml:"clonefrom"`
	} `yaml:"tags"`
}

var defaultPatroniSpec *PatroniV11Specification

func BuildPatroniSpec(clusterSpec *ClusterSpecification, hostDiscoverySpec *HostDiscoverySpecification) (patroniSpec *PatroniV11Specification, err error) {
	patroniSpec, err = DefaultPatroniSpec()
	if err != nil {
		return
	}
	patroniSpec.MergeClusterSpec(clusterSpec, hostDiscoverySpec)
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

func (patroniSpec *PatroniV11Specification) MergeClusterSpec(clusterSpec *ClusterSpecification, hostDiscoverySpec *HostDiscoverySpecification) {
	appuserName := clusterSpec.Postgresql.Appuser.Username
	replicationUsername := appuserName
	patroniSpec.Etcd.Host = clusterSpec.Etcd.URI
	patroniSpec.Scope = clusterSpec.Cluster.Scope
	patroniSpec.Name = clusterSpec.Cluster.Name
	patroniSpec.Bootstrap.PgHba = []string{
		fmt.Sprintf("host replication %s 0.0.0.0/0 md5", replicationUsername),
		"host postgres all 0.0.0.0/0 md5",
	}
	patroniSpec.Bootstrap.Users.Postgres.Password = clusterSpec.Postgresql.Admin.Password
	patroniSpec.Postgresql.Authentication.Replication.Username = clusterSpec.Postgresql.Appuser.Username
	patroniSpec.Postgresql.Authentication.Replication.Password = clusterSpec.Postgresql.Appuser.Password
	patroniSpec.Postgresql.Authentication.Superuser.Username = clusterSpec.Postgresql.Superuser.Username
	patroniSpec.Postgresql.Authentication.Superuser.Password = clusterSpec.Postgresql.Superuser.Password

	patroniSpec.Postgresql.ConnectAddress = fmt.Sprintf("%s:%s", hostDiscoverySpec.IP, hostDiscoverySpec.Port5432)
	patroniSpec.Restapi.ConnectAddress = fmt.Sprintf("%s:%s", hostDiscoverySpec.IP, hostDiscoverySpec.Port8008)
}

func (patroniSpec *PatroniV11Specification) String() string {
	bytes, err := yaml.Marshal(patroniSpec)
	if err != nil {
		panic(err)
	}
	return string(bytes[:])
}

func (patroniSpec *PatroniV11Specification) CreateConfigFile(path string) (err error) {
	data, err := yaml.Marshal(patroniSpec)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, data, 0644)
	return
}
