/*
 * Created on Mon Dec 23 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package join

import (
	"fmt"

	"github.com/AndrewSerra/crowdsourced-testcases/internal/api"
	datastorage "github.com/AndrewSerra/crowdsourced-testcases/internal/data-storage"
	"github.com/spf13/cobra"
)

// joinCmd represents the join command
var JoinCmd = &cobra.Command{
	Use:   "join",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		join_tk, _ := cmd.Flags().GetString("token")
		courseid, _ := cmd.Flags().GetInt("courseid")

		if join_tk == "" {
			fmt.Println("join token is required")
			return
		}

		if courseid == -1 {
			fmt.Println("courseid is required")
			return
		}

		profile, err := datastorage.GetActiveUserProfile()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = api.AcceptStudentForCourse(courseid, profile.Id, join_tk)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Successfully joined course")
		// TODO: Store locally
	},
}

func init() {
	JoinCmd.Flags().IntP("courseid", "c", -1, "Course id to join")
	JoinCmd.Flags().StringP("token", "t", "", "Token to join course")

	JoinCmd.MarkFlagRequired("courseid")
	JoinCmd.MarkFlagRequired("token")
}
