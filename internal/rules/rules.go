package rules

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/karthikkashyap98/sweeperd/internal/utils"
)

type Rule struct {
	Name    string  `yaml:"rule"`
	Enabled bool    `yaml:"enabled"`
	Match   Match   `yaml:"match"`
	Action  Action  `yaml:"action"`
	Options Options `yaml:"options"`
}

type Match struct {
	Folder        string   `yaml:"folder"`
	extensions    []string `yaml:"extensions"`
	OlderThanDays int      `yaml:"older_than_days"`
}

type ActionType int

const (
	Move ActionType = iota
	Delete
	Rename
)

type Action struct {
	Type   ActionType `yaml:"type"`
	Target string     `yaml:"target"`
}

type Options struct {
	DryRun bool `yaml:"dry_run"`
	Log    bool `yaml:"log"`
}

func LoadRules(rulesFile string) *Rule {
	var rule Rule

	data, err := os.ReadFile(rulesFile)
	fmt.Println(rulesFile)
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &rule)
	if err != nil {
		utils.PrintError(err)
		os.Exit(1)
	}

	return &rule
}
