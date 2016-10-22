package clicmd

import (
	"testing"

	"github.com/dingotiles/dingo-postgresql-agent/testutil"
)

func TestTestAPI_filterWaleEnvVars(t *testing.T) {
	t.Parallel()
	if len(filterWaleEnvVarsFromList([]string{})) != 0 {
		t.Fatalf("Empty list input should return empty list")
	}

	result := filterWaleEnvVarsFromList([]string{"INVALID", "FOO_IGNORE=1", "WAL_S3_BUCKET=1", "WALE_KEEP=1", "WABS_KEEP=1", "GOOGLE_KEEP=1", "IGNORE=1", "SWIFT_KEEP=1"})
	if !testutil.TestEqStringArray(result, []string{"WAL_S3_BUCKET=1", "WALE_KEEP=1", "WABS_KEEP=1", "GOOGLE_KEEP=1", "SWIFT_KEEP=1"}) {
		t.Fatalf("IGNORE/KEEP result should return 5 KEEP items, returned: %#v", result)
	}

	result = filterWaleEnvVarsFromList([]string{"WALE_KEEP=1", "WALE_IGNORE="})
	if !testutil.TestEqStringArray(result, []string{"WALE_KEEP=1"}) {
		t.Fatalf("Should ignore env vars without values, returned: %#v", result)
	}
}
