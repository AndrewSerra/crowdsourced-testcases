/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package cmd

import (
	"os"

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
