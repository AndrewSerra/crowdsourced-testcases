/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// Student
func CreateStudent(student NewStudent) (int, error) {
	db := GetDB()

	result, err := db.Exec(
		"INSERT INTO students (first_name, last_name, email) VALUES (?, ?, ?)",
		student.FirstName, student.LastName, student.Email)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return int(insertedId), nil
}

func CreateStudentBatch(students []NewStudent) error {
	db := GetDB()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, student := range students {
		_, err := tx.Exec("INSERT INTO students (first_name, last_name, email) VALUES (?, ?, ?)",
			student.FirstName, student.LastName, student.Email)

		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func GetStudent(id string) (*Student, error) {
	db := GetDB()

	var student Student
	row := db.QueryRow("SELECT id, first_name, last_name, email FROM students WHERE id = ?", id)
	err := row.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &student, nil
}

func GetAnonymousStudent(id string) (*StudentAnonymous, error) {
	db := GetDB()

	var student StudentAnonymous
	row := db.QueryRow("SELECT id, anon_name FROM students WHERE id = ?", id)
	err := row.Scan(&student.Id, &student.AnonymousId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &student, nil
}

// Instructor
func CreateInstructor(instructor NewInstructor) (int, error) {
	db := GetDB()

	result, err := db.Exec(
		"INSERT INTO instructors (first_name, last_name, email) VALUES (?, ?, ?)",
		instructor.FirstName, instructor.LastName, instructor.Email)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return int(insertedId), nil
}

// Course
func CreateCourse(course NewCourse) (int, error) {
	db := GetDB()

	result, err := db.Exec("INSERT INTO courses (title, owner_id) VALUES (?, ?)", course.Name, course.OwnerId)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return int(insertedId), nil
}

func GetCourse(id int) (*Course, error) {
	db := GetDB()

	var course Course
	row := db.QueryRow("SELECT id, title, owner_id FROM courses WHERE id = ?", id)
	err := row.Scan(&course.Id, &course.Name, &course.OwnerId)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &course, nil
}

func DeleteCourse(courseid int) (int, error) {
	db := GetDB()

	result, err := db.Exec("DELETE FROM courses WHERE id = ?", courseid)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return int(rowCount), nil
}

// Assignment
func CreateAssignment(assignment NewAssignment) (int, error) {
	db := GetDB()

	startDate, err := time.Parse(time.RFC3339, assignment.StartDate)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	endDate, err := time.Parse(time.RFC3339, assignment.EndDate)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	result, err := db.Exec(
		"INSERT INTO assignments (title, course_id, start_date, end_date) VALUES (?, ?, ?, ?)",
		assignment.Name, assignment.CourseId, startDate, endDate)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return int(insertedId), nil
}

func GetAssignment(courseid int, assignmentid int) (*Assignment, error) {
	db := GetDB()

	var assignment Assignment
	row := db.QueryRow(
		"SELECT id, title, course_id, start_date, end_date, is_open, is_published FROM assignments WHERE id = ? AND course_id = ?", assignmentid, courseid)
	err := row.Scan(&assignment.Id, &assignment.Name, &assignment.CourseId, &assignment.StartDate, &assignment.EndDate, &assignment.IsOpen, &assignment.IsPublished)

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &assignment, nil
}

func GetAssignmentsForCourse(courseid int) ([]*Assignment, error) {
	db := GetDB()

	rows, err := db.Query("SELECT id, title, course_id, start_date, end_date, is_open, is_published FROM assignments WHERE course_id = ?", courseid)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var assignments []*Assignment = []*Assignment{}
	for rows.Next() {
		var assignment Assignment
		err := rows.Scan(&assignment.Id, &assignment.Name, &assignment.CourseId, &assignment.StartDate, &assignment.EndDate, &assignment.IsOpen, &assignment.IsPublished)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}

	return assignments, nil
}

func DeleteAssignment(courseid int, assignmentid int) (int, error) {
	db := GetDB()

	result, err := db.Exec("DELETE FROM assignments WHERE id = ? AND course_id = ?", assignmentid, courseid)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return int(rowCount), nil
}

func updateAssignment(cid int, aid int, val int, col string) (int, error) {
	db := GetDB()

	if val < 0 || val > 1 {
		return -1, fmt.Errorf("invalid value for %s", col)
	}

	if col != "is_open" && col != "is_published" {
		return -1, errors.New("invalid column name")
	}

	query := fmt.Sprintf("UPDATE assignments SET %s = ? WHERE id = ? AND course_id = ?", col)
	result, err := db.Exec(query, val, aid, cid)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return int(rowCount), nil
}

func publishingAssignmentGrades(courseid int, assignmentid int, value int) (int, error) {
	return updateAssignment(courseid, assignmentid, value, "is_published")
}

func openingAssignmentGrades(courseid int, assignmentid int, value int) (int, error) {
	return updateAssignment(courseid, assignmentid, value, "is_open")
}

func SetPublishedAssignment(courseid int, assignmentid int) (int, error) {
	return publishingAssignmentGrades(courseid, assignmentid, 1)
}

func ClearPublishedAssignment(courseid int, assignmentid int) (int, error) {
	return publishingAssignmentGrades(courseid, assignmentid, 0)
}

func SetOpenAssignment(courseid int, assignmentid int) (int, error) {
	return openingAssignmentGrades(courseid, assignmentid, 1)
}

func ClearOpenAssignment(courseid int, assignmentid int) (int, error) {
	return openingAssignmentGrades(courseid, assignmentid, 0)
}

// Submission
func createSubmission() {

}

func CreateTestCaseSubmission() {

}

func CreateAssignmentSubmission() {

}
