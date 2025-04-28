package interfaces

import (
	"backend-go/internal/dto"
)

type MailingListRepository interface {
	Save(newMailingList *dto.MailingList) error
}
