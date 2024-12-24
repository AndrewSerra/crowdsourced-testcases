/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package main

import (
	"os"

	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/edit"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/join"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/new"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/profile"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/publish"
	"github.com/AndrewSerra/crowdsourced-testcases/cmd/testsource-cli/stats"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crowdsourced-testcases",
	Short: "Crowdsourcing test cases for programming assignments",
	Long: `Crowdsourcing test cases for programming assignments. Students can submit
	test cases to be run for all submissions in a course assignment. Instructors can
	create courses and assignment for a roster of students.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(edit.EditCmd)
	rootCmd.AddCommand(join.JoinCmd)
	rootCmd.AddCommand(new.NewCmd)
	rootCmd.AddCommand(profile.ProfileCmd)
	rootCmd.AddCommand(publish.PublishCmd)
	rootCmd.AddCommand(stats.StatsCmd)
}
