package config

// PatroniSpecification is constructed based on ClusterSpecification provided by the API.
// It is converted to a patroni.yml and used by Patroni to configure & run PostgreSQL.
type PatroniSpecification struct {
	Scope    string `yaml:"scope"`
	TTL      int    `yaml:"ttl"`
	LoopWait int    `yaml:"loop_wait"`
	Etcd     struct {
		Host  string `yaml:"host"`
		Scope string `yaml:"scope"`
		TTL   int    `yaml:"ttl"`
	} `yaml:"etcd"`
	Postgresql struct {
		Admin struct {
			Password string `yaml:"password"`
			Username string `yaml:"username"`
		} `yaml:"admin"`
		ConnectAddress       string   `yaml:"connect_address"`
		CreateReplicaMethod  []string `yaml:"create_replica_method"`
		DataDir              string   `yaml:"data_dir"`
		Listen               string   `yaml:"listen"`
		MaximumLagOnFailover int      `yaml:"maximum_lag_on_failover"`
		Name                 string   `yaml:"name"`
		Parameters           struct {
			ArchiveCommand      string `yaml:"archive_command"`
			ArchiveMode         string `yaml:"archive_mode"`
			ArchiveTimeout      string `yaml:"archive_timeout"`
			HotStandby          string `yaml:"hot_standby"`
			ListenAddresses     string `yaml:"listen_addresses"`
			MaxConnections      int    `yaml:"max_connections"`
			MaxReplicationSlots int    `yaml:"max_replication_slots"`
			MaxWalSenders       int    `yaml:"max_wal_senders"`
			MaxWalSize          string `yaml:"max_wal_size"`
			MinWalSize          string `yaml:"min_wal_size"`
			Port                int    `yaml:"port"`
			WalKeepSegments     int    `yaml:"wal_keep_segments"`
			WalLevel            string `yaml:"wal_level"`
			WalLogHints         string `yaml:"wal_log_hints"`
		} `yaml:"parameters"`
		PgHba        []string `yaml:"pg_hba"`
		Pgpass       string   `yaml:"pgpass"`
		RecoveryConf struct {
			RestoreCommand string `yaml:"restore_command"`
		} `yaml:"recovery_conf"`
		Replication struct {
			Network  string `yaml:"network"`
			Password string `yaml:"password"`
			Username string `yaml:"username"`
		} `yaml:"replication"`
		Scope     string `yaml:"scope"`
		Superuser struct {
			Password string `yaml:"password"`
			Username string `yaml:"username"`
		} `yaml:"superuser"`
		UseSlots bool `yaml:"use_slots"`
		WalE     struct {
			Command                       string `yaml:"command"`
			Envdir                        string `yaml:"envdir"`
			NoMaster                      int    `yaml:"no_master"`
			Retries                       int    `yaml:"retries"`
			ThresholdBackupSizePercentage int    `yaml:"threshold_backup_size_percentage"`
			ThresholdMegabytes            int    `yaml:"threshold_megabytes"`
			UseIam                        int    `yaml:"use_iam"`
		} `yaml:"wal_e"`
	} `yaml:"postgresql"`
	RestAPI struct {
		ConnectAddress string `yaml:"connect_address"`
		Listen         string `yaml:"listen"`
	} `yaml:"restapi"`
}
