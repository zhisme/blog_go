package interfaces

import (
	"backend-go/internal/dto"
)

type MailingListValidator interface {
	Validate(mailingList *dto.MailingList) error
}
