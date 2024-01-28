package cl30_test

import (
	"errors"
	"testing"

	cl "github.com/opencl-go/cl30"
)

func allPlatforms(tb testing.TB) []cl.PlatformID {
	tb.Helper()
	ids, err := cl.PlatformIDs()
	if err != nil {
		if errors.Is(err, cl.StatusError(-1001)) {
			tb.Errorf("failed to query platform IDs: %v", err)
		}
		return nil
	}
	return ids
}

func TestPlatforms(t *testing.T) {
	platformIDs := allPlatforms(t)
	if len(platformIDs) == 0 {
		t.Skipf("no platforms available")
	}
	for _, platformID := range platformIDs {
		name, err := cl.PlatformInfoString(platformID, cl.PlatformNameInfo)
		if err != nil {
			t.Logf("failed to query name of platform: %v", err)
		}
		t.Logf("Platform <%s>\n", name)
	}
}
