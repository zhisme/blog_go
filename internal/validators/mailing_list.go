package validators

import (
	"backend-go/internal/dto"
	"errors"
	"regexp"
)

type MailingListValidator struct{}

func NewMailingListValidator() *MailingListValidator {
	return &MailingListValidator{}
}

func (m *MailingListValidator) Validate(mailingList *dto.MailingList) error {
	if err := m.validateEmail(mailingList.Email); err != nil {
		return err
	}

	if err := m.validateUsername(mailingList.Email); err != nil {
		return err
	}

	return nil
}

func (m *MailingListValidator) validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (m *MailingListValidator) validateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	return nil
}
