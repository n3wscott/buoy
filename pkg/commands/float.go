package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"tableflip.dev/buoy/pkg/float"
)

func addFloatCmd(root *cobra.Command) {
	var domain string
	var release string
	var strict bool

	var floatCmd = &cobra.Command{
		Use:   "float go.mod",
		Short: "Find latest versions of dependencies based on a release.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gomod := args[0]

			refs, err := float.Float(gomod, release, domain, strict)
			if err != nil {
				return err
			}

			for _, r := range refs {
				if r != "" {
					fmt.Printf("%s\n", r)
				}
			}
			return nil
		},
	}

	floatCmd.Flags().StringVarP(&domain, "domain", "d", "knative.dev", "domain filter")
	floatCmd.Flags().StringVarP(&release, "release", "r", "", "release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]")
	_ = floatCmd.MarkFlagRequired("release")
	floatCmd.Flags().BoolVarP(&strict, "strict", "s", false, "strict - only select and return tagged modules")

	root.AddCommand(floatCmd)
}
