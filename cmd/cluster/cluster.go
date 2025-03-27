/*
Copyright © 2025 Phillip Nguyen <png9981@gmail.com>
*/
package cmd

import (
	clusterDestroy "github.com/saldyy/kube-cli/cmd/cluster/destroy"
	clusterInit "github.com/saldyy/kube-cli/cmd/cluster/init"
	clusterResume "github.com/saldyy/kube-cli/cmd/cluster/resume"
	clusterUpdate "github.com/saldyy/kube-cli/cmd/cluster/update"

	"github.com/spf13/cobra"
)

// ClusterCmd represents the cluster command
var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	ClusterCmd.AddCommand(clusterInit.InitCmd)
	ClusterCmd.AddCommand(clusterUpdate.UpdateCmd)
	ClusterCmd.AddCommand(clusterDestroy.DestroyCmd)
	ClusterCmd.AddCommand(clusterResume.ResumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
