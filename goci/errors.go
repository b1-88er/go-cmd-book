package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation error")
)

// type step int

// const (
// 	build step = iota
// )

// func (s step) String() string {
// 	switch s {
// 	case build:
// 		return "go build"
// 	default:
// 		return "unknown"
// 	}
// }

type stepErr struct {
	step  string
	msg   string
	cause error
}

func (e *stepErr) Error() string {
	return fmt.Sprintf("step: %s, msg: %s, cause: %v", e.step, e.msg, e.cause)
}

func (e *stepErr) Is(target error) bool {
	if t, ok := target.(*stepErr); ok {
		return t.step == e.step
	}
	return false
}

func (e *stepErr) Unwrap() error {
	return e.cause
}
