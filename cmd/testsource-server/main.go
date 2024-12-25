/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import "github.com/gin-gonic/gin"

func main() {
	defer SafelyCloseDB()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	courseGroup := router.Group("/courses")
	{
		courseGroup.POST("/", CreateCourseHandler)
		courseGroup.GET("/:cid", GetCourseHandler)
		courseGroup.DELETE("/:cid", DeleteCourseHandler)
		courseGroup.POST("/:cid/roster", CreateRosterHandler)

		assignmentGroup := courseGroup.Group("/:cid/assignments")
		{
			assignmentGroup.POST("/", CreateAssignmentHandler)
			assignmentGroup.GET("/", GetAssignmentsForCourseHandler)
			assignmentGroup.GET("/:aid", GetAssignmentHandler)
			assignmentGroup.DELETE("/:aid", DeleteAssignmentHandler)
			assignmentGroup.POST("/:aid/*action", AssignmentActionsHandler)
		}
	}

	instructorGroup := router.Group("/instructors")
	{
		instructorGroup.POST("/", CreateInstructorHandler)
	}

	router.Run(":8080")
}
