package minecraft

import (
	"path/filepath"
	"testing"
)

func TestResolveConfiguredServerJar(t *testing.T) {
	binaryDirectory := t.TempDir()
	externalJar := filepath.Join(t.TempDir(), "custom.jar")
	for _, tc := range []struct{ configured, want string }{
		{"", filepath.Join(binaryDirectory, "server.jar")},
		{"paper-26.3.jar", filepath.Join(binaryDirectory, "paper-26.3.jar")},
		{filepath.Join("server", "custom.jar"), filepath.Join(binaryDirectory, "server", "custom.jar")},
		{externalJar, externalJar},
	} {
		if got := ResolveServerJarPath(tc.configured, binaryDirectory); got != tc.want {
			t.Fatalf("%q: got %q want %q", tc.configured, got, tc.want)
		}
	}
}
