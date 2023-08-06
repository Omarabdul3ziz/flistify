package cmds

import "github.com/spf13/cobra"

var RootCommand = &cobra.Command{
	Use:   "flistify",
	Short: "Tool to build, run and store flists.",
}
