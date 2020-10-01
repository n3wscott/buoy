package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"tableflip.dev/buoy/pkg/needs"
)

func addNeedsCmd(root *cobra.Command) {
	var domain string
	var dot bool

	var floatCmd = &cobra.Command{
		Use:   "needs go.mod",
		Short: "Find dependencies based on a base import domain.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gomod := args[0]

			if dot {
				graph, err := needs.Dot(gomod, domain)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", graph)
				return nil
			}

			packages, err := needs.Needs(gomod, domain)
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

	floatCmd.Flags().StringVarP(&domain, "domain", "d", "knative.dev", "domain filter")
	floatCmd.Flags().BoolVar(&dot, "dot", false, "Produce a .dot file output for use with Graphviz.")

	root.AddCommand(floatCmd)
}
