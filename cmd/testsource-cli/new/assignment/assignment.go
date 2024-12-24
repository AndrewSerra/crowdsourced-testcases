/*
 * Created on Sun Dec 22 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package assignment

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/AndrewSerra/crowdsourced-testcases/internal/api"
	"github.com/spf13/cobra"
)

type datetimeInput [2]string

// assignmentCmd represents the assignment command
var AssignmentCmd = &cobra.Command{
	Use:   "assignment",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assignment called")
		title := args[0]
		courseId, _ := cmd.Flags().GetString("course")

		var dates datetimeInput = [2]string{}

		startDate := readInputFromStdIn("Start Date (yyyy-mm-dd): ", "^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
		startTime := readInputFromStdIn("Start Time (hh:mm): ", "^[0-9]{2}:[0-9]{2}$")
		dates[0] = fmt.Sprintf("%s %s:00", startDate, startTime)

		endDate := readInputFromStdIn("End Date (yyyy-mm-dd): ", "^[0-9]{4}-[0-9]{2}-[0-9]{2}$")
		endTime := readInputFromStdIn("End Time (hh:mm): ", "^[0-9]{2}:[0-9]{2}$")
		dates[1] = fmt.Sprintf("%s %s:00", endDate, endTime)

		valid, err := verifyDatetimeInput(dates)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !valid {
			fmt.Println("Invalid datetime input")
			os.Exit(1)
		}

		assignmentId, err := api.CreateAssignmentForCourse(title, courseId, struct{ Start, End string }{
			Start: dates[0],
			End:   dates[1],
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("New Asssingment Id: %d Name: %s\n", assignmentId, title)
		// TODO: create assignment in storage
	},
}

func verifyDatetimeInput(datetimeIn datetimeInput) (bool, error) {
	timestamp1, err := time.Parse(time.DateTime, datetimeIn[0])
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	timestamp2, err := time.Parse(time.DateTime, datetimeIn[1])
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	if time.Now().After(timestamp1) {
		return false, errors.New("start date must be in the future")
	}

	if time.Now().After(timestamp2) {
		return false, errors.New("end date must be in the future")
	}

	if timestamp1.After(timestamp2) {
		return false, errors.New("end date must be after start date")
	}

	return true, nil
}

func readInputFromStdIn(prompt string, pattern string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> %s", prompt)
		switch line, err := reader.ReadString('\n'); err {
		case nil:
			line := strings.TrimSpace(line)
			match, err := regexp.MatchString(pattern, line)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if match {
				return line
			} else {
				fmt.Println("Invalid input")
			}
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Printf("Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
	}
}

func init() {
	AssignmentCmd.Flags().String("course", "", "Course id of the assignment")

	AssignmentCmd.MarkFlagRequired("course")
}
