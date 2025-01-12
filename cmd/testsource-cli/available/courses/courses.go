/*
 * Created on Sun Jan 12 2025
 *
 * Copyright © 2025 Andrew Serra <andy@serra.us>
 */
package courses

import (
	"fmt"

	datastorage "github.com/AndrewSerra/crowdsourced-testcases/internal/data-storage"
	"github.com/spf13/cobra"
)

// coursesCmd represents the courses command
var CoursesCmd = &cobra.Command{
	Use:   "courses",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
 and usage of using your command. For example:

 Cobra is a CLI library for Go that empowers applications.
 This application is a tool to generate the needed files
 to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		courses, err := datastorage.GetAvailableCoursesForActiveProfile()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(courses) == 0 {
			fmt.Println("No available courses")
			return
		}

		for _, c := range courses {
			fmt.Println(c.Name, "-", len(c.Assignments), "assignments")
		}
	},
}
