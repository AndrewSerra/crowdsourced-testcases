/*
 * Created on Tue Dec 24 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package use

import (
	"fmt"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var UseCmd = &cobra.Command{
	Use:   "use",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
 and usage of using your command. For example:

 Cobra is a CLI library for Go that empowers applications.
 This application is a tool to generate the needed files
 to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use called")
	},
}
