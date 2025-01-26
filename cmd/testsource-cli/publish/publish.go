/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package publish

import (
	"fmt"

	api "github.com/AndrewSerra/crowdsourced-testcases/internal/api"
	datastorage "github.com/AndrewSerra/crowdsourced-testcases/internal/data-storage"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var PublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		courseName, _ := cmd.Flags().GetString("course-name")
		assignmentName, _ := cmd.Flags().GetString("assignment-name")
		unpublish, _ := cmd.Flags().GetBool("unpublish")

		if assignmentName == "" || courseName == "" {
			fmt.Println("No assignment name or course name provided")
			return
		}

		profile, err := datastorage.GetActiveUserProfile()
		if err != nil {
			fmt.Println(err)
			return
		}

		ownerid := profile.Id

		courses, err := datastorage.GetAvailableCoursesForActiveProfile()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, course := range courses {
			if courseName == course.Name {
				assignment := course.GetAssignmentByName(assignmentName)

				if assignment != nil {
					if unpublish {
						err = api.UnpublishAssignmentGradesByName(ownerid, course.Id, assignment.Id)
						if err != nil {
							fmt.Println(err)
							return
						}

						fmt.Println("Unublished assignment grades successfully")
					} else {
						err = api.PublishAssignmentGradesByName(ownerid, course.Id, assignment.Id)
						if err != nil {
							fmt.Println(err)
							return
						}

						fmt.Println("Published assignment grades successfully")
					}
					return
				}
			}
		}
		fmt.Println("Assignment name or course name not found")
	},
}

func init() {
	PublishCmd.Flags().StringP("course-name", "c", "", "Course name of assignment")
	PublishCmd.Flags().StringP("assignment-name", "a", "", "Assignment name to publish")

	PublishCmd.Flags().Bool("unpublish", false, "Unpublish the assignment")

	PublishCmd.MarkFlagRequired("course-name")
	PublishCmd.MarkFlagRequired("assignment-name")

}
