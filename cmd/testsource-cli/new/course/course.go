/*
 * Created on Sun Dec 22 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package course

import (
	"fmt"

	"github.com/spf13/cobra"
)

// courseCmd represents the course command
var CourseCmd = &cobra.Command{
	Use:   "course",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("course called")
	},
}

func init() {

}
