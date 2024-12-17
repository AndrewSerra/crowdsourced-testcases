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

// Assignment
func CreateAssignmentHandler(c *gin.Context) {
	var courseUri struct {
		CourseId string `uri:"cid"`
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
					"error": fmt.Sprintf("course '%s' does not exist", assignment.CourseId),
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

	err = DeleteAssignment(cid, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func PublishAssignmentGradesHandler(c *gin.Context) {
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

	err = PublishAssignmentGrades(cid, aid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
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

func CreateRosterHandler(c *gin.Context) {
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

	err = CreateStudentBatch(students)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
