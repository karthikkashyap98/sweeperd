package cli

import (
	"context"
	"fmt"

	"github.com/karthikkashyap98/sweeperd/internal/config"
	"github.com/karthikkashyap98/sweeperd/internal/executor/actions"
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
	rule := rules.LoadRules(rulesPath)

	fmt.Println(cfg)
	fmt.Println(rule)
	return nil
}

func runHandler(cmd *cobra.Command, args []string) error {

	return nil
}

func BuildAndRun(ctx context.Context, r rules.Rule) error {
	matcher := BuildMatcherFromRule(r)

	act, err := actions.NewAction(r.Action, r.Match.Folder, r.Action.Target, matcher)
	if err != nil {
		return err
	}

	files, err := act.Plan(ctx)
	if err != nil {
		return err
	}

	return act.Execute(ctx, files)
}
