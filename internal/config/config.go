package config

type Config struct {
	RulesPath    string   `yaml:"rules_path"`
	LogPath      string   `yaml:"log_path"`
	TrashPath    string   `yaml:"trash_path"`
	WatchEnabled bool     `yaml:"watch_enabled"`
	WatchFolders []string `yaml:"watch_folders"`
	DebounceMS   int      `yaml:"debounce_ms"`
	Parallelism  int      `yaml:"parallelism"`
	DryRunGlobal bool     `yaml:"dry_run_global"`
}
