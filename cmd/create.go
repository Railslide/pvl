package cmd

import (
	"github.com/railslide/pvl/internal"
	"github.com/spf13/cobra"
)

var path string

// CreateCmd represents the add command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Pyright configuration file",
	Long:  `Create Pyright configuration file in the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.CreateConfig(path)

	},
}

func init() {
	CreateCmd.Flags().StringVar(&path, "path", "", "custom path to the virtualenv")

	rootCmd.AddCommand(CreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
