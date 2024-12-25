/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package edit

import (
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/edit/roster"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: nil,
}

func init() {
	EditCmd.AddCommand(roster.RosterCmd)
}
