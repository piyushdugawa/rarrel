package cmd

import (
	"github.com/piyushdugawa/rarrel/internal/lua"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <plugin>",
	Args:  cobra.ExactArgs(1),
	Short: "Create Project Using Plugin mentioned after",
	Run: func(cmd *cobra.Command, args []string) {
		plugin := args[0]
		lua.LoadLua(plugin)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
