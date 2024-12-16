/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

type SubmissionStatus int

const (
	SUBMISSION_SUBMITTED SubmissionStatus = iota
	SUBMISSION_SUCCESS
	SUBMISSION_FAILED
)

// Base structs
type submission struct {
	CourseId    string `json:"course_id"`
	OwnerId     string `json:"owner_id"`
	SubmittedAt string `json:"submitted_at"`
}

type person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Derived structs and models
type Course struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Assignment struct {
	Id string `json:"id"`
}

type TestCaseSubmission struct {
	Id string `json:"id"`
	submission
}

type AssignmentSubmission struct {
	Id          string           `json:"id"`
	Status      SubmissionStatus `json:"status"`
	IsPublished bool             `json:"is_published"`
	submission
}

type Student struct {
	Id string `json:"id"`
	person
}

type StudentAnonymous struct {
	Id          string `json:"id"`
	AnonymousId string `json:"anonymous_id"`
	person
}

type Instructor struct {
	Id string `json:"id"`
	person
}
