package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"tableflip.dev/buoy/pkg/check"
)

func addCheckCmd(root *cobra.Command) {
	var domain string
	var release string
	var verbose bool

	var cmd = &cobra.Command{
		Use:   "check go.mod",
		Short: "Determine if this module has release branches or releases available from each dependency for a given release.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gomod := args[0]

			err := check.Check(gomod, release, domain, verbose)
			if errors.Is(err, check.DependencyErr) {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
				return nil
			}

			return err
		},
	}

	cmd.Flags().StringVarP(&domain, "domain", "d", "knative.dev", "domain filter")
	cmd.Flags().StringVar(&release, "release", "r", "release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]")
	_ = cmd.MarkFlagRequired("release")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output.")

	root.AddCommand(cmd)
}
