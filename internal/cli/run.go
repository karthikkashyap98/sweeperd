package cli

import (
	"fmt"

	"github.com/karthikkashyap98/sweeperd/internal/config"
	"github.com/karthikkashyap98/sweeperd/internal/rules"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var ruleName string
	var matchFolder string
	var concurrency int

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Run all enabled rules once and exit",
		PersistentPreRunE: preRunHandler,
		RunE:              runHandler,
	}

	cmd.Flags().StringVar(&ruleName, "rule", "", "Run a single rule by name")
	cmd.Flags().StringVar(&matchFolder, "match-folder", "", "Override rule's match.folder at runtime")
	cmd.Flags().IntVar(&concurrency, "concurrency", 0, "Limit concurrent operations (default: from config)")

	return cmd
}

func preRunHandler(cmd *cobra.Command, args []string) error {
	cfg := config.LoadConfig(cfgPath)
	rules := rules.LoadRules(rulesPath)

	fmt.Println(cfg)
	fmt.Println(rules)
	return nil
}

func runHandler(cmd *cobra.Command, args []string) error {
	return nil
}
