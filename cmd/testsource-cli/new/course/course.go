/*
 * Created on Sun Dec 22 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package course

import (
	"fmt"

	"github.com/AndrewSerra/crowdsourced-testcases/internal/api"
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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("course called")
		title := args[0]
		ownerId, _ := cmd.Flags().GetString("owner")

		courseId, err := api.CreateCourseForInstructor(title, ownerId)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("New Course Id: %d Name: %s\n", courseId, title)
		// TODO: create course in storage
	},
}

func init() {
	CourseCmd.Flags().String("owner", "", "Instructor id of the course")
	CourseCmd.MarkFlagRequired("owner")
}
