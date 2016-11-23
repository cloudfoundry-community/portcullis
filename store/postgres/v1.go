package postgres

import "fmt"

type v1 struct {
}

func (v v1) migrate() error {
	//TODO
	return fmt.Errorf("Not yet implemented")
}

func (v v1) version() int {
	return 1
}
