package version

// Version information set at build time by goreleaser via ldflags.
// Use: -ldflags "-X github.com/borland502/wordgen/internal/version.Version=v1.2.3 ..."
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)
