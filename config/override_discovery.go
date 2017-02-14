package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// OverrideSpecification collects all user overrides from environment variables
type OverrideSpecification struct {
	PostgresRecovery *PostgresRecoveryOverrideSpecification
}

// PostgresRecoveryOverrideSpecification imports postgres recovery overrides from environment variables
type PostgresRecoveryOverrideSpecification struct {
	TargetName      string `envconfig:"pg_recovery_target_name"`
	TargetTime      string `envconfig:"pg_recovery_target_time"`
	TargetXid       string `envconfig:"pg_recovery_target_xid"`
	TargetInclusive bool   `envconfig:"pg_recovery_target_inclusive"`
	TargetTimeline  string `envconfig:"pg_recovery_target_timeline"`
	TargetAction    string `envconfig:"pg_recovery_target_action"`
}

// OverrideSpec collects all the user overrides into single struct
func OverrideSpec() (override *OverrideSpecification) {
	override = &OverrideSpecification{}
	override.PostgresRecovery = loadPostgresRecoveryOverrides()
	return
}

func loadPostgresRecoveryOverrides() (spec *PostgresRecoveryOverrideSpecification) {
	spec = &PostgresRecoveryOverrideSpecification{}
	err := envconfig.Process("pg_recovery", spec)
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}
