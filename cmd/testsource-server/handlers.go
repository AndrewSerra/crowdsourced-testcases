/*
 * Created on Mon Dec 16 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

const (
	TOKEN_HEADER = "X-TestSource-Token"
)

// Assignment
func CreateAssignmentHandler(c *gin.Context) {
	var courseUri struct {
		CourseId int `uri:"cid"`
	}
	var assignment NewAssignment

	if err := c.ShouldBindUri(&courseUri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&assignment); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "request body is empty",
			})
		} else if validationErrors, ok := err.(validator.ValidationErrors); ok {
			missingFields := GetMissingFieldString(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("missing required fields: %s", missingFields),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	assignment.CourseId = courseUri.CourseId

	assignment_id, err := CreateAssignment(assignment)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == FOREIGN_KEY_NO_EXIST_ERROR {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("course '%d' does not exist", assignment.CourseId),
				})
			}
			if mysqlError.Number == DUPLICATE_ENTRY_ERROR {
				c.JSON(http.StatusConflict, gin.H{
					"error": fmt.Sprintf("assignment '%s' already exists", assignment.Name),
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(200, gin.H{
		"id": assignment_id,
	})
}

func GetAssignmentHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")
	param_aid := c.Params.ByName("aid")

	if param_cid == "" || param_aid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	aid, err := strconv.Atoi(param_aid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	assignment, err := GetAssignment(cid, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else if assignment == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "assignment not found",
		})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func GetAssignmentsForCourseHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")

	if param_cid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	assignments, err := GetAssignmentsForCourse(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

func DeleteAssignmentHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")
	param_aid := c.Params.ByName("aid")

	if param_cid == "" || param_aid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	aid, err := strconv.Atoi(param_aid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rowCount, err := DeleteAssignment(cid, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else if rowCount == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.Status(http.StatusOK)
}

func AssignmentActionsHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")
	param_aid := c.Params.ByName("aid")

	if param_cid == "" || param_aid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	aid, err := strconv.Atoi(param_aid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	action := c.Param("action")[1:]
	var rowCount int

	switch action {
	case "publish":
		rowCount, err = SetPublishedAssignment(cid, aid)
	case "unpublish":
		rowCount, err = ClearPublishedAssignment(cid, aid)
	case "open":
		rowCount, err = SetOpenAssignment(cid, aid)
	case "close":
		rowCount, err = ClearOpenAssignment(cid, aid)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("unknown action '%s'", action),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else if rowCount == 0 {
		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusOK)
	}
}

// Instructor
func CreateInstructorHandler(c *gin.Context) {
	var instructor NewInstructor

	if err := c.ShouldBindJSON(&instructor); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "request body is empty",
			})
		} else if validationErrors, ok := err.(validator.ValidationErrors); ok {
			missingFields := GetMissingFieldString(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("missing required fields: %s", missingFields),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	instructor_id, err := CreateInstructor(instructor)
	if err != nil {
		if mysqlError, ok := err.(*mysql.MySQLError); ok {
			if mysqlError.Number == DUPLICATE_ENTRY_ERROR {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("instructor '%s' already exists", instructor.Email),
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// TODO: send email to instructor

	c.JSON(200, gin.H{
		"id": instructor_id,
	})
}

func GetInstructorHandler(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	instructor, err := GetInstructorByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else if instructor == nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "instructor not found",
		})
		return
	}
}

// Course
func CreateCourseHandler(c *gin.Context) {
	var course NewCourse

	if err := c.ShouldBindJSON(&course); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "request body is empty",
			})
		} else if validationErrors, ok := err.(validator.ValidationErrors); ok {
			missingFields := GetMissingFieldString(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("missing required fields: %s", missingFields),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	course_id, err := CreateCourse(course)
	if err != nil {
		if dbError, ok := err.(*mysql.MySQLError); ok {
			if dbError.Number == FOREIGN_KEY_NO_EXIST_ERROR {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("instructor '%s' does not exist", course.OwnerId),
				})
			}
			if dbError.Number == DUPLICATE_ENTRY_ERROR {
				c.JSON(http.StatusConflict, gin.H{
					"error": fmt.Sprintf("course '%s' already exists", course.Name),
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(200, gin.H{
		"id": course_id,
	})
}

func GetCourseHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")

	if param_cid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	course, err := GetCourse(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else if course == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "course not found",
		})
		return
	}

	c.JSON(http.StatusOK, course)
}

func DeleteCourseHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")

	if param_cid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	rowCount, err := DeleteCourse(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else if rowCount == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.Status(http.StatusOK)
}

func CreateRosterHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")

	if param_cid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var students []NewStudent

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &students)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = CreateStudentBatch(cid, students)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func ApproveCourseStudentRegistrationHandler(c *gin.Context) {
	param_cid := c.Params.ByName("cid")
	param_sid := c.Params.ByName("sid")

	if param_cid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	cid, err := strconv.Atoi(param_cid)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if param_sid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	sid, err := strconv.Atoi(param_sid)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var code struct {
		EntryCode string `json:"entry_code" binding:"required"`
	}

	err = json.Unmarshal(body, &code)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// err = CompleteRegisterationToCourse(cid, token, code.EntryCode)
	isFullfilled, err := CompleteRegisterationToCourse(cid, sid, code.EntryCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	if !isFullfilled {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

// Student
func GetStudentHandler(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	student, err := GetStudentByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else if student == nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "student not found",
		})
		return
	} else {
		c.JSON(http.StatusOK, student)
	}
}

func VerifyStudentHandler(c *gin.Context) {
	param_sid := c.Params.ByName("sid")

	if param_sid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	sid, err := strconv.Atoi(param_sid)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	isVerified, err := VerifyStudentEmail(sid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else if !isVerified {
		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusOK)
	}
}
