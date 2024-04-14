package storage

import "errors"

var (
	ErrEntityExists   = errors.New("entity already exists")
	ErrEntityNotFound = errors.New("entity not found")
)
