package library_repo

import "errors"

var (
	ErrTitleNotFound = errors.New("No books found with that title")
)
