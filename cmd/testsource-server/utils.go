/*
 * Created on Tue Dec 17 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package main

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetMissingFieldString(validationErrors validator.ValidationErrors) string {
	missing := []string{}

	for _, validationError := range validationErrors {
		if validationError.Tag() == "required" {
			missing = append(missing, validationError.Field())
		}
	}

	return strings.Join(missing, ", ")
}
