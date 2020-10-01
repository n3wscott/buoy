package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"tableflip.dev/buoy/pkg/needs"
)

func addNeedsCmd(root *cobra.Command) {
	var domain string
	var dot bool

	var cmd = &cobra.Command{
		Use:   "needs go.mod",
		Short: "Find dependencies based on a base import domain.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gomods := args

			if dot {
				graph, err := needs.Dot(gomods, domain)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", graph)
				return nil
			}

			_, packages, err := needs.Needs(gomods, domain)
			if err != nil {
				return err
			}

			for _, p := range packages {
				if p != "" {
					fmt.Printf("%s\n", p)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&domain, "domain", "d", "knative.dev", "domain filter")
	cmd.Flags().BoolVar(&dot, "dot", false, "Produce a .dot file output for use with Graphviz.")

	root.AddCommand(cmd)
}
