/*
 * Created on Tue Dec 17 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package profile

import (
	"fmt"

	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var ProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manages profiles for the CLI operations",
	Long:  `Different emails can be used as different profiles for the CLI operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profile called")

	},
}
