package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/karthikkashyap98/sweeperd/internal/config"
	"github.com/karthikkashyap98/sweeperd/internal/executor/actions"
	"github.com/karthikkashyap98/sweeperd/internal/rules"
	"github.com/spf13/cobra"
)

type ctxKey string

const (
	cfgKey  ctxKey = "config"
	ruleKey ctxKey = "rules"
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

	ctx := context.WithValue(cmd.Context(), cfgKey, cfg)
	ctx = context.WithValue(ctx, ruleKey, rule)

	cmd.SetContext(ctx)

	return nil
}

func runHandler(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	ruleset, ok := ctx.Value(ruleKey).(*rules.Rule)
	if !ok {
		return errors.New("rules not found in context")
	}

	BuildAndRun(context.Background(), *ruleset)
	return nil
}

type MatcherFunc func(path string, d fs.DirEntry) bool

func (m MatcherFunc) Match(path string, d fs.DirEntry) bool {
	return m(path, d)
}

var (
	MatchNothing MatcherFunc = func(path string, d fs.DirEntry) bool { return false }
)

func BuildMatcherFromRule(r rules.Rule) actions.Matcher {
	// TODO: This should not happen for every file in the directory
	if !r.Enabled {
		return MatchNothing
	}

	extensions := r.Match.Extensions

	if len(extensions) > 0 {
		return MatcherFunc(func(path string, d fs.DirEntry) bool {
			if d.IsDir() {
				return false
			}

			ext := strings.TrimPrefix(
				strings.ToLower(filepath.Ext(d.Name())),
				".",
			)

			for _, e := range extensions {
				if ext == strings.ToLower(e) {
					return true
				}
			}
			return false
		})
	}

	return MatcherFunc(func(path string, d fs.DirEntry) bool {
		return false
	})
}

func BuildAndRun(ctx context.Context, r rules.Rule) error {
	if r.Enabled == false {
		fmt.Println("Rule is disabled")
		return nil
	}

	matcher := BuildMatcherFromRule(r)

	act, err := actions.NewAction(r, matcher)
	if err != nil {
		return err
	}

	files, err := act.Plan(ctx)
	if err != nil {
		return err
	}

	return act.Execute(ctx, files)
}
