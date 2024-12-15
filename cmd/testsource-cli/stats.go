/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stats called")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
