package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"tableflip.dev/buoy/pkg/git"

	"github.com/spf13/cobra"
	"tableflip.dev/buoy/pkg/check"
)

func addCheckCmd(root *cobra.Command) {
	var domain string
	var release string
	var rulesetFlag string
	var ruleset git.RulesetType
	var verbose bool

	var cmd = &cobra.Command{
		Use:   "check go.mod",
		Short: "Determine if this module has a ref for each dependency for a given release based on a ruleset.",
		Long: `
The check command is used to evaluate if each dependency for the given module
meets the requirements for cutting a release branch. If the requirements are
met based on the ruleset selected, the command will exit with code 0, otherwise
an error message is generated and the with the failed dependencies and exit
code 1. Errors are written to stderr. Verbose output is written to stdout.

Rulesets,
  Release          check requires all dependencies to have tagged releases.
  Branch           check requires all dependencies to have a release branch.
  ReleaseOrBranch  check will use rule (Release || Branch).

`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Validation
			ruleset = git.Ruleset(rulesetFlag)
			if ruleset == git.InvalidRule {
				return fmt.Errorf("invalid ruleset, please select one of: [%s]", strings.Join(git.Rulesets(), ", "))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			gomod := args[0]

			err := check.Check(gomod, release, domain, ruleset, verbose)
			if errors.Is(err, check.DependencyErr) {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
				return nil
			}

			return err
		},
	}

	cmd.Flags().StringVarP(&domain, "domain", "d", "knative.dev", "domain filter")
	cmd.Flags().StringVarP(&release, "release", "r", "", "release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]")
	_ = cmd.MarkFlagRequired("release")
	cmd.Flags().StringVar(&rulesetFlag, "ruleset", git.ReleaseOrReleaseBranchRule.String(), fmt.Sprintf("The ruleset to evaluate the dependency refs. Rulesets: [%s]", strings.Join(git.Rulesets(), ", ")))
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output.")

	root.AddCommand(cmd)
}
