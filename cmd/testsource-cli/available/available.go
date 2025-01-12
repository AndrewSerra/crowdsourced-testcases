/*
 * Created on Sun Jan 12 2025
 *
 * Copyright © 2025 Andrew Serra <andy@serra.us>
 */

package available

import (
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/available/assignments"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/available/courses"
	"github.com/spf13/cobra"
)

// availableCmd represents the available command
var AvailableCmd = &cobra.Command{
	Use:   "available",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
 and usage of using your command. For example:

 Cobra is a CLI library for Go that empowers applications.
 This application is a tool to generate the needed files
 to quickly create a Cobra application.`,
	Run: nil,
}

func init() {
	AvailableCmd.AddCommand(courses.CoursesCmd)
	AvailableCmd.AddCommand(assignments.AssignmentsCmd)
}
