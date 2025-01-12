/*
 * Created on Sun Jan 12 2025
 *
 * Copyright © 2025 Andrew Serra <andy@serra.us>
 */
package assignments

import (
	"fmt"

	datastorage "github.com/AndrewSerra/crowdsourced-testcases/internal/data-storage"
	"github.com/spf13/cobra"
)

// assignmentsCmd represents the assignments command
var AssignmentsCmd = &cobra.Command{
	Use:   "assignments",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
  and usage of using your command. For example:

  Cobra is a CLI library for Go that empowers applications.
  This application is a tool to generate the needed files
  to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		courseId, _ := cmd.Flags().GetInt("course-id")

		var assignments []datastorage.Assignment
		var err error

		if courseId == -1 {
			assignments, err = datastorage.GetAllAssignmentsForActiveProfile()
		} else {
			assignments, err = datastorage.GetAvailableAssignmentsForCourse(courseId)
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		if len(assignments) == 0 {
			fmt.Println("No assignments found")
			return
		}

		for _, assignment := range assignments {
			fmt.Println(assignment)
		}
	},
}

func init() {
	AssignmentsCmd.Flags().IntP("course-id", "c", -1, "Course-id to filter assignments by")
}
