package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation error")
)

type step int

const (
	build step = iota
)

func (s step) String() string {
	switch s {
	case build:
		return "go build"
	default:
		return "unknown"
	}
}

type stepErr struct {
	step  step
	msg   string
	cause error
}

func (e *stepErr) Error() string {
	return fmt.Sprintf("step: %s, msg: %s, cause: %v", e.step, e.msg, e.cause)
}

func (e *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}
	return t.step == e.step
}

func (e *stepErr) Unwrap() error {
	return e.cause
}
