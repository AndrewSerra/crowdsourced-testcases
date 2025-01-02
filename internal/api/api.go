/*
 * Created on Sun Dec 22 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const url = "http://localhost:8080"

type createResponseBody struct {
	Id int `json:"id" binding:"required"`
}

type persontype string

const (
	INSTRUCTOR persontype = "instructor"
	STUDENT    persontype = "student"
)

func CreateCourseForInstructor(title string, ownerId string) (int, error) {
	data := map[string]interface{}{
		"name":     title,
		"owner_id": ownerId,
	}
	body, err := generateRequestBodyJSON(data)
	if err != nil {
		return -1, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/courses", url), "application/json", body)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var jsonRespBody createResponseBody
	err = json.Unmarshal(respBody, &jsonRespBody)

	if err != nil {
		return -1, err
	}

	return jsonRespBody.Id, nil
}

func CreateAssignmentForCourse(title string, courseId string, duration struct{ Start, End string }) (int, error) {
	startdate, err := time.Parse(time.DateTime, duration.Start)
	if err != nil {
		return -1, err
	}
	enddate, err := time.Parse(time.DateTime, duration.End)
	if err != nil {
		return -1, err
	}

	data := map[string]interface{}{
		"name":       title,
		"course_id":  courseId,
		"start_date": startdate,
		"end_date":   enddate,
	}
	body, err := generateRequestBodyJSON(data)
	if err != nil {
		return -1, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/courses/%s/assignments", url, courseId), "application/json", body)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var jsonRespBody createResponseBody
	err = json.Unmarshal(respBody, &jsonRespBody)

	if err != nil {
		return -1, err
	}

	return jsonRespBody.Id, nil
}

func AcceptStudentForCourse(courseId int, studentId int, token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}

	if courseId == -1 {
		return fmt.Errorf("courseId is required")
	}

	if studentId == -1 {
		return fmt.Errorf("studentId is required")
	}

	body, err := generateRequestBodyJSON(map[string]string{
		"entry_code": token,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/courses/%d/students/%d/accept", url, courseId, studentId), "application/json", body)

	// request, err := http.NewRequest("POST", fmt.Sprintf("%s/courses/%d/students/%d/accept", url, courseId, studentId), nil)
	// if err != nil {
	// 	return err
	// }
	// request.Header.Set("X-TestSource-Token", token)

	// resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("email is not verified for account")
	}
	return nil
}

func CreateCourseStudentRoster(courseid int, students []NewStudent) error {
	if len(students) == 0 {
		return fmt.Errorf("at least one student is required")
	}

	body, err := generateRequestBodyJSON(students)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("%s/courses/%d/roster", url, courseid), "application/json", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println(res)
		return fmt.Errorf("failed to create roster")
	}

	return nil
}

func GetInstructorByEmail(email string) (*Person, error) {
	return getPersonByEmail(email, INSTRUCTOR)
}

func GetStudentByEmail(email string) (*Person, error) {
	return getPersonByEmail(email, STUDENT)
}

func getPersonByEmail(email string, t persontype) (*Person, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%ss?email=%s", url, t, email))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	} else if resp.StatusCode == http.StatusOK {
		var person Person
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &person)
		if err != nil {
			return nil, err
		}
		return &person, nil
	} else {
		return nil, fmt.Errorf("unexpected error: %s", resp.Status)
	}
}

func generateRequestBodyJSON(data interface{}) (*bytes.Buffer, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	newbody := bytes.NewBuffer(jsonBytes)

	return newbody, nil
}
