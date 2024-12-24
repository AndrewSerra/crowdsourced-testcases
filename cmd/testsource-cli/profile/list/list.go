/*
 * Created on Tue Dec 24 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package list

import (
	"fmt"
	"os"
	"strings"

	datastorage "github.com/AndrewSerra/crowdsourced-testcases/internal/data-storage"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
 and usage of using your command. For example:

 Cobra is a CLI library for Go that empowers applications.
 This application is a tool to generate the needed files
 to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := datastorage.GetUserProfileList()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		activeName := datastorage.GetActiveProfileName()
		headerStr := "Available Profiles:"
		headerLen := len(headerStr)

		if len(profiles) == 0 {
			fmt.Println("No profiles found")
			return
		}

		fmt.Println(headerStr)
		fmt.Println(strings.Repeat("-", headerLen))

		for _, profile := range profiles {
			if profile == activeName {
				fmt.Printf("  * %s (active)\n", profile)
				continue
			}
			fmt.Printf("  * %s\n", profile)
		}
	},
}
