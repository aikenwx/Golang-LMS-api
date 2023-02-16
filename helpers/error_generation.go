package helpers

import "fmt"

func GenerateInvalidRequestsError() error {
	return fmt.Errorf("Required fields for the request are missing or have invalid types")
}
