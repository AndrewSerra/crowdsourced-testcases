/*
 * Created on Tue Dec 24 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package new

import (
	"fmt"

	"github.com/AndrewSerra/crowdsourced-testcases/internal/api"
	datastorage "github.com/AndrewSerra/crowdsourced-testcases/internal/data-storage"
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
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		// token, _ := cmd.Flags().GetString("token")

		isUsed, err := datastorage.IsEmailUsedInProfile(email)

		if err != nil {
			fmt.Println(err)
			return
		}

		if isUsed {
			fmt.Println("Email already in use")
			return
		}

		isInstructor, _ := cmd.Flags().GetBool("instructor")
		var target api.Person

		if isInstructor {
			instructor, err := api.GetInstructorByEmail(email)
			if err != nil {
				fmt.Println(err)
				return
			}

			if instructor == nil {
				fmt.Println("Instructor not found")
				return
			}

			target = *instructor
		} else {
			student, err := api.GetStudentByEmail(email)
			if err != nil {
				fmt.Println(err)
				return
			}

			if student == nil {
				fmt.Println("Student not found")
				return
			}

			target = *student
		}

		err = datastorage.CreateNewUserProfile(name, datastorage.UserProfile(target))

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Profile '%s' created\n", name)
	},
}

func init() {
	NewCmd.Flags().StringP("name", "n", "", "Name of the profile")
	NewCmd.Flags().StringP("email", "e", "", "Email of the profile")
	NewCmd.Flags().StringP("token", "t", "", "Token of the profile")
	NewCmd.Flags().Bool("instructor", false, "Treat profile as instructor")

	NewCmd.MarkFlagRequired("name")
	NewCmd.MarkFlagRequired("email")
	NewCmd.MarkFlagRequired("token")
}
