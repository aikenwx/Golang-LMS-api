package helpers

import (
	"fmt"
	"net/mail"
	"regexp"
)

// ValidateEmailFormat custom email validation handler
func ValidateEmailFormat(email string) error {
	_, err := mail.ParseAddress(email)

	if err != nil {
		return fmt.Errorf(`The email address %s has an invalid format`, email)
	}

	return nil
}

func ValidateEmailAddresses(emails []string) error {
	for _, email := range emails {
		err := ValidateEmailFormat(email)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindValidEmailsInText(text string) []string {
	regexPattern := regexp.MustCompile(`@(?i)\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}\b`)

	emails := regexPattern.FindAllString(text, -1)

	// remove @ sign
	emails = Map(emails, func(email string) string {
		return email[1:]
	})

	// return only valid emails
	emails = Filter(emails, func(email string) bool {
		return ValidateEmailFormat(email) == nil
	})

	return emails
}
