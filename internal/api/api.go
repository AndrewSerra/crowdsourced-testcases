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

func AcceptStudentForCourse(courseId string, token string) error {
	if token == "" {
		return fmt.Errorf("token is required")
	}

	if courseId == "" {
		return fmt.Errorf("courseId is required")
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/courses/%s/accept", url, courseId), nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-TestSource-Token", token)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: Do error handling

	return nil
}

func generateRequestBodyJSON(data interface{}) (*bytes.Buffer, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	newbody := bytes.NewBuffer(jsonBytes)

	return newbody, nil
}
