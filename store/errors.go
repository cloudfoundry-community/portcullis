package store

import "fmt"

//ErrNotFound is the error that should be returned if a Get or Delete
// request is performed on a non-existent mapping
var ErrNotFound = fmt.Errorf("The requested mapping was not found in the store")

//ErrDuplicate is the error that should be returned if an attempt to add a
// mapping to the store where a mapping with that name already exists
var ErrDuplicate = fmt.Errorf("The given mapping already exists in the store")
