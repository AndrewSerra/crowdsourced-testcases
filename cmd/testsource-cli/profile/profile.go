/*
 * Created on Tue Dec 17 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package profile

import (
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/profile/delete"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/profile/list"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/profile/new"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/profile/use"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var ProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manages profiles for the CLI operations",
	Long:  `Different emails can be used as different profiles for the CLI operations.`,
	Run:   nil,
}

func init() {
	ProfileCmd.AddCommand(delete.DeleteCmd)
	ProfileCmd.AddCommand(new.NewCmd)
	ProfileCmd.AddCommand(use.UseCmd)
	ProfileCmd.AddCommand(list.ListCmd)
}
