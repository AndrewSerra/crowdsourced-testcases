/*
 * Created on Sun Dec 22 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package assignment

import (
	"fmt"

	"github.com/spf13/cobra"
)

// assignmentCmd represents the assignment command
var AssignmentCmd = &cobra.Command{
	Use:   "assignment",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assignment called")
	},
}

func init() {

}
