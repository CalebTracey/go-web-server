package config

import "fmt"

func MissingField(field string) error {
	return fmt.Errorf("%v not found in config file", field)
}
