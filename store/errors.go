package store

import "fmt"

//TODO: Change all the errors over to struct types

//Error is what errors originating from the Store package should be
type Error struct {
	message string
	code    int
}

func (e Error) Error() string {
	return e.message
}

const (
	codeNotFound = iota
	codeDuplicate
	codeInvalid
)

//ErrNotFound is the error that should be returned if a Get or Delete
// request is performed on a non-existent mapping
var ErrNotFound = fmt.Errorf("The requested mapping was not found in the store")

//ErrDuplicate is the error that should be returned if an attempt to add a
// mapping to the store where a mapping with that name already exists
var ErrDuplicate = fmt.Errorf("The given mapping already exists in the store")

//NewErrInvalid makes an error of the type that should be returned if there is
// something about a mapping which violates a value constraint (e.g. length, type)
func NewErrInvalid(mess string) Error {
	return Error{message: mess, code: codeInvalid}
}

//IsErrInvalid returns true if this is a store.Error signaling that a value was
// not valid. False otherwise.
func IsErrInvalid(e error) bool {
	castErr, isStoreError := e.(Error)
	if !isStoreError || castErr.code != codeInvalid {
		return false
	}
	return true
}
