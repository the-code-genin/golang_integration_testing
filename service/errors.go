package service

import "errors"

var (
	ErrInternal     = errors.New("an internal error occurred")
	ErrNoteNotFound = errors.New("no note was found with the ID")
)
