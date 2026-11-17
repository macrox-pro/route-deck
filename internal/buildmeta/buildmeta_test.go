package buildmeta

import "testing"

func TestDisplayVersion(t *testing.T) {
	tests := []struct {
		name string
		ver  string
		want string
	}{
		{name: "default_dev", ver: "dev", want: "dev"},
		{name: "empty", ver: "", want: "dev"},
		{name: "semver", ver: "1.2.3", want: "v1.2.3"},
		{name: "already_prefixed", ver: "v0.1.0", want: "v0.1.0"},
		{name: "whitespace", ver: " 2.0.0 ", want: "v2.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev := Version
			t.Cleanup(func() { Version = prev })
			Version = tt.ver
			if got := DisplayVersion(); got != tt.want {
				t.Fatalf("DisplayVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPlistVersion(t *testing.T) {
	tests := []struct {
		name string
		ver  string
		want string
	}{
		{name: "dev", ver: "dev", want: "0.0.0"},
		{name: "empty", ver: "", want: "0.0.0"},
		{name: "semver", ver: "1.2.3", want: "1.2.3"},
		{name: "strip_v", ver: "v0.1.0", want: "0.1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev := Version
			t.Cleanup(func() { Version = prev })
			Version = tt.ver
			if got := PlistVersion(); got != tt.want {
				t.Fatalf("PlistVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}
