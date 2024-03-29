package main

import (
	"errors"
)

var (
	ErrNotNumber        = errors.New("not a number")
	ErrInvalidColum     = errors.New("invalid column")
	ErrNoFiles          = errors.New("no files to process")
	ErrInvalidOperation = errors.New("invalid operation")
)
