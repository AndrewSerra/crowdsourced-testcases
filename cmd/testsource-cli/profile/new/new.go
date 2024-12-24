/*
 * Created on Tue Dec 24 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package new

import (
	"fmt"

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
		fmt.Println("new called")
		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")

		var profile map[string]string = map[string]string{
			"name":  name,
			"email": email,
		}

		fmt.Printf("profile: %+v\n", profile)
		// TODO:
		// 1 - check if profile exists (email)
		// 2 - create profile if not exist
	},
}

func init() {
	NewCmd.Flags().StringP("name", "n", "", "Name of the profile")
	NewCmd.Flags().StringP("email", "e", "", "Email of the profile")

	NewCmd.MarkFlagRequired("name")
	NewCmd.MarkFlagRequired("email")
}
