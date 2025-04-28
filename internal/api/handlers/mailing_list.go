package handlers

import (
  "time"
	"backend-go/internal/dto"
	"backend-go/internal/validators"
	"backend-go/internal/repositories"
)

func HandleCreate(newMailingList dto.MailingList) (error, dto.MailingList) {
  validator := validators.NewMailingListValidator()
  if err := validator.Validate(&newMailingList); err != nil {
    return err, newMailingList
  }

  mailingList := &dto.MailingList{
    Username: newMailingList.Username,
    Email:   newMailingList.Email,
    CreatedAt: time.Now(),
  }

	repo := repositories.NewCsvMailingListRepository("mailing_list.csv")

  repo.Save(mailingList)

  return nil, *mailingList
}
