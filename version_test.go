package cl30_test

import (
	"testing"

	cl "github.com/opencl-go/cl30"
)

func TestVersionOf(t *testing.T) {
	t.Parallel()
	type args struct {
		major int
		minor int
		patch int
	}
	tt := []struct {
		name string
		args args
		want cl.Version
	}{
		{name: "0.0.0", args: args{major: 0, minor: 0, patch: 0}, want: 0},
		{name: "1.0.0", args: args{major: 1, minor: 0, patch: 0}, want: 0x00400000},
		{name: "1.1.1", args: args{major: 1, minor: 1, patch: 1}, want: 0x00401001},
		{name: "4.5.6", args: args{major: 4, minor: 5, patch: 6}, want: 0x01005006},
		{name: "0.0.-1", args: args{major: 0, minor: 0, patch: -1}, want: 0x00000FFF},
		{name: "0.-1.0", args: args{major: 0, minor: -1, patch: 0}, want: 0x003FF000},
		{name: "-1.0.0", args: args{major: -1, minor: 0, patch: 0}, want: 0xFFC00000},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := cl.VersionOf(tc.args.major, tc.args.minor, tc.args.patch); got != tc.want {
				t.Errorf("VersionOf() = 0x%08X, want 0x%08X", got, tc.want)
			}
		})
	}
}

func TestVersionComponents(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name  string
		ver   cl.Version
		major int
		minor int
		patch int
	}{
		{name: "VersionMin", ver: cl.VersionMin, major: 0, minor: 0, patch: 0},
		{name: "VersionMax", ver: cl.VersionMax, major: 0x3FF, minor: 0x3FF, patch: 0xFFF},
		{name: "1.2.3", ver: cl.VersionOf(1, 2, 3), major: 1, minor: 2, patch: 3},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.ver.Major(); got != tc.major {
				t.Errorf("Major() = 0x%03X, want 0x%03X", got, tc.major)
			}
			if got := tc.ver.Minor(); got != tc.minor {
				t.Errorf("Minor() = 0x%03X, want 0x%03X", got, tc.minor)
			}
			if got := tc.ver.Patch(); got != tc.patch {
				t.Errorf("Patch() = 0x%03X, want 0x%03X", got, tc.patch)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	t.Parallel()
	tt := []struct {
		ver  cl.Version
		want string
	}{
		{ver: cl.VersionOf(0, 0, 0), want: "0.0.0"},
		{ver: cl.VersionOf(1, 2, 3), want: "1.2.3"},
		{ver: cl.VersionOf(10, 0, 100), want: "10.0.100"},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()
			if got := tc.ver.String(); got != tc.want {
				t.Errorf("String() = '%s', want '%s'", got, tc.want)
			}
		})
	}
}
