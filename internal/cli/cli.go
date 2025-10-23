package cli

import (
	"fmt"
	"os"

	"github.com/karthikkashyap98/sweeperd/internal/utils"
	"github.com/karthikkashyap98/sweeperd/pkg/constants"
	"github.com/spf13/cobra"
)

var (
	cfgPath   string
	rulesPath string
	logDir    string
	dryRun    bool
	jsonOut   bool
	quiet     bool
	verbose   bool
)

func ValidateFiles(cmd *cobra.Command, args []string) error {
	if !utils.FileExists(cfgPath) {
		err := fmt.Errorf("Config File does not exist %s", cfgPath)
		return err
	}

	if !utils.FileExists(rulesPath) {
		err := fmt.Errorf("Rules not found in the directory %s", rulesPath)
		return err
	}

	return nil
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "sweeperd",
		Short:             "Rule-based file organiser and daemon",
		PersistentPreRunE: ValidateFiles,
	}

	cmd.PersistentFlags().StringVar(&cfgPath, "config", constants.DefaultConfigFile, "Path to config file")
	cmd.PersistentFlags().StringVar(&rulesPath, "rules", constants.DefaultRulesDir, "Directory containing rule files")
	cmd.PersistentFlags().StringVar(&logDir, "log", constants.DefaultLogsDir, "Directory for logs")
	cmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Force dry-run regardless of config/rules")
	cmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Emit JSON output")
	cmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "Suppress non-error output")
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Increase verbosity")

	cmd.AddCommand(newRunCmd())

	return cmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
