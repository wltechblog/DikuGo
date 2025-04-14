package command

import (
	"errors"
)

// Command errors
var (
	ErrCommandNotFound   = errors.New("command not found")
	ErrWrongPosition     = errors.New("you are in the wrong position for that")
	ErrInsufficientLevel = errors.New("you are not high enough level for that")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrNotImplemented    = errors.New("command not implemented yet")
)
