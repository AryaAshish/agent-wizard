package buildinfo

var (
	// Version is set at link time by scripts/release/build_assets.sh (-ldflags -X ...).
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
