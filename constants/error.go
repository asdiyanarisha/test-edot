package constants

import "errors"

var (
	ErrorPostAlreadyInserted = errors.New("post already inserted")
	ErrorPostNotFound        = errors.New("post not found")
)
