/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package new

import (
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/new/assignment"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/new/course"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: nil,
}

func init() {
	NewCmd.AddCommand(assignment.AssignmentCmd)
	NewCmd.AddCommand(course.CourseCmd)
}
