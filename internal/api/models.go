/*
 * Created on Wed Dec 25 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */
package api

type NewStudent struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
