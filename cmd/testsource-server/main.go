/*
 * Created on Sat Dec 14 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	router.Run(":8080")
}
