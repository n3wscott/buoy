package commands

import "github.com/spf13/cobra"

func New() *cobra.Command {
	var buoyCmd = &cobra.Command{
		Use:   "buoy",
		Short: "Introspect go module dependencies.",
	}

	addFloatCmd(buoyCmd)
	addNeedsCmd(buoyCmd)

	return buoyCmd
}
