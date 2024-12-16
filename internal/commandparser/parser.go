/*
 * Created on Sun Dec 15 2024
 *
 * Copyright © 2024 Andrew Serra <andy@serra.us>
 */

package commandparser

type ActionName int
type ActionPrompt struct{}

const (
	ADD_COURSE ActionName = iota
	ADD_ASSIGNMENT
	DELETE_COURSE
	DELETE_ASSIGNMENT
	UPDATE_START_DATE
	UPDATE_END_DATE
	UPDATE_ASSIGNMENT_AVAIL
)

func New() *ActionPrompt {
	return &ActionPrompt{}
}

type Action[T any] struct {
	Name    ActionName
	Old     T
	Replace T
}
