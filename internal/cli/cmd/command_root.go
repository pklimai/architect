package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "architect",
	Short: "A genearator for architect applications.",
	Long:  "Aarchitect is a CLI tool for fast generating uniform GO microservices, based on gitlab.com/zigal0/architect",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init initialize and bind all cli commands
func init() {
	// root subcommmands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(showCmd)

	// generate sub commands
	generateCmd.AddCommand(servicesCmd)

	// add sub commands
	addCmd.AddCommand(managerCmd)
	addCmd.AddCommand(subManagerCmd)
	addCmd.AddCommand(repositoryCmd)
	addCmd.AddCommand(protoServiceCmd)
	addCmd.AddCommand(clientCmd)

	// show
	showCmd.AddCommand(architectureCmd)
}
