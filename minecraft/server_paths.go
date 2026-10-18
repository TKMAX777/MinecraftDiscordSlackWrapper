package minecraft

import "path/filepath"

// ResolveServerJarPath resolves the configured jar relative to the wrapper binary.
// Only an empty setting uses Minecraft's conventional server.jar name.
func ResolveServerJarPath(configured, binaryDirectory string) string {
	if configured == "" {
		configured = "server.jar"
	}
	if !filepath.IsAbs(configured) {
		configured = filepath.Join(binaryDirectory, configured)
	}
	return filepath.Clean(configured)
}
