/*
 * Created on Tue Dec 24 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package roster

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/AndrewSerra/crowdsourced-testcases/internal/api"
	"github.com/spf13/cobra"
)

var headers_required = [...]string{
	"first_name",
	"last_name",
	"email",
}

// rosterCmd represents the roster command
var RosterCmd = &cobra.Command{
	Use:   "roster",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cid, _ := cmd.Flags().GetInt("courseid")
		filename, _ := cmd.Flags().GetString("file")
		hasNoHeaders, _ := cmd.Flags().GetBool("noHeader")

		if cid == 0 {
			fmt.Println("No course id provided")
			return
		}

		if filename == "" {
			fmt.Println("No file provided")
			return
		}

		if _, err := os.Stat(filename); err != nil {
			fmt.Println(err)
			return
		}

		ext := path.Ext(filename)

		if ext != ".csv" {
			fmt.Println("file must be a csv")
			return
		}

		data, err := readCsv(filename, hasNoHeaders)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = api.CreateCourseStudentRoster(cid, data)

		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	RosterCmd.Flags().IntP("courseid", "c", 0, "Course id to join")
	RosterCmd.Flags().StringP("file", "f", "", "File to read roster from")
	RosterCmd.Flags().Bool("noHeader", false, "File has no header")

	RosterCmd.MarkFlagRequired("courseid")
	RosterCmd.MarkFlagRequired("file")
}

func readCsv(filename string, hasNoHeaders bool) ([]api.NewStudent, error) {
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(records[0]) != len(headers_required) {
		return nil, fmt.Errorf("csv must have exactly 3 headers, in order: first_name, last_name, email")
	}

	if hasNoHeaders {
		return getRosterStudentList(records), nil
	} else {
		valid, err := checkCsvHeaders(records[0])
		if !valid {
			return nil, fmt.Errorf("invalid headers: %s", err)
		}
		return getRosterStudentList(records[1:]), nil
	}
}

func checkCsvHeaders(headers []string) (bool, error) {

	if len(headers) > 3 {
		return false, fmt.Errorf("csv must have exactly 3 headers: first_name, last_name, email")
	}

	missing := make([]string, 0, 3)
	found := func() map[string]bool {
		temp := map[string]bool{}
		for _, header := range headers_required {
			temp[header] = false
		}
		return temp
	}()

	for _, header := range headers {
		if _, ok := found[header]; ok {
			found[header] = true
		}
	}

	for header := range found {
		if !found[header] {
			missing = append(missing, header)
		}
	}

	if len(missing) > 0 {
		return false, fmt.Errorf("missing headers: %s", strings.Join(missing, ", "))
	}

	return true, nil
}

func getRosterStudentList(records [][]string) []api.NewStudent {
	students := make([]api.NewStudent, 0, len(records))
	for _, record := range records {
		students = append(students, api.NewStudent{
			FirstName: record[0],
			LastName:  record[1],
			Email:     record[2],
		})
	}
	return students
}
