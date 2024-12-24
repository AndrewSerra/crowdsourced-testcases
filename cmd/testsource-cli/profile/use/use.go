/*
 * Created on Tue Dec 24 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package use

import (
	"fmt"

	datastorage "github.com/AndrewSerra/crowdsourced-testcases/internal/data-storage"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var UseCmd = &cobra.Command{
	Use:   "use",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
 and usage of using your command. For example:

 Cobra is a CLI library for Go that empowers applications.
 This application is a tool to generate the needed files
 to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		currentProfile := datastorage.GetActiveProfileName()
		err := datastorage.SetNewActiveProfileState(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Switched from %s -> %s\n", currentProfile, name)
	},
}
