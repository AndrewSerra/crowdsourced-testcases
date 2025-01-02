/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import "time"

type SubmissionStatus int

const (
	SUBMISSION_SUBMITTED SubmissionStatus = iota
	SUBMISSION_SUCCESS
	SUBMISSION_FAILED
)

// Base structs
type submission struct {
	CourseId     int    `uri:"cid" json:"course_id" binding:"required"`
	OwnerId      int    `json:"owner_id" binding:"required"`
	AssignmentId string `uri:"aid" json:"assignment_id" binding:"required"`
	// SubmittedAt string `json:"submitted_at"`     // Maybe add it later for offline submission?
}

type person struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

// Derived structs and models
type NewCourse struct {
	Name    string `json:"name" binding:"required"`
	OwnerId string `json:"owner_id" binding:"required"`
}
type Course struct {
	Id    int    `json:"id"`
	Token string `json:"join_tk"`
	NewCourse
}

type NewAssignment struct {
	CourseId  int    `json:"course_id"`
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type Assignment struct {
	Id          int       `json:"id"`
	IsOpen      bool      `json:"is_open"`
	IsPublished bool      `json:"is_published"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	NewAssignment
}

type NewTestCaseSubmission struct {
	InputData      string `json:"input_data" binding:"required"`
	ExpectedResult string `json:"expected_result" binding:"required"`
	submission
}

type TestCaseSubmission struct {
	Id int `json:"id"`
	NewTestCaseSubmission
}

type NewAssignmentSubmission submission

type AssignmentSubmission struct {
	Id     int              `json:"id"`
	Status SubmissionStatus `json:"grading_status"`
	NewAssignmentSubmission
}

type NewStudent person

type Student struct {
	Id int `json:"id"`
	NewStudent
}

type StudentAnonymous struct {
	Id          int    `json:"id"`
	AnonymousId string `json:"anonymous_id"`
}

type NewInstructor person

type Instructor struct {
	Id int `json:"id"`
	NewInstructor
}
