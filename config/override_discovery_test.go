package config

import (
	"fmt"
	"os"
	"testing"
)

func TestOverrideDiscovery(t *testing.T) {
	os.Clearenv()
	os.Setenv("PG_RECOVERY_TARGET_TIMELINE", "latest")

	spec := OverrideSpec()
	if spec.PostgresRecovery.TargetTimeline != "latest" {
		fmt.Printf("%#v\n", *spec.PostgresRecovery)
		t.Fatalf("TargetTimeline should be 'latest', got '%s'", spec.PostgresRecovery.TargetTimeline)
	}
}
