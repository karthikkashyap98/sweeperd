package config

import (
	"github.com/goccy/go-yaml"
	"github.com/karthikkashyap98/sweeperd/internal/utils"
	"os"
)

type Config struct {
	RulesPath    string `yaml:"rules_path"`
	LogPath      string `yaml:"log_path"`
	TrashPath    string `yaml:"trash_path"`
	WatchEnabled bool   `yaml:"watch_enabled"`
	DebounceMS   int    `yaml:"debounce_ms"`
	Parallelism  int    `yaml:"parallelism"`
}

func LoadConfig(configFilePath string) *Config {
	var cfg Config

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}

	return &cfg
}
