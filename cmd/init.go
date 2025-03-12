/*
Copyright © 2025 Phillip Nguyen <png9981@gmail.com>
*/
package cmd

import (
	"github.com/saldyy/kube-cli/services"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "The init command sets up a local Kubernetes cluster using Minikube.",
	Long: `The init command sets up a local Kubernetes cluster using Minikube.
It ensures that Minikube is installed, starts a new cluster with the specified
configuration, and verifies that essential components like the Kubernetes API
server and the default namespace are running. This command simplifies the setup process,
making it easy to start developing and testing Kubernetes applications locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		services.InitCluster()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
