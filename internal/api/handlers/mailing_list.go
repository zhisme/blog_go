package handlers

import (
	"backend-go/internal/dto"
	"backend-go/internal/repositories"
	"backend-go/internal/validators"
	"time"
)

func HandleCreate(newMailingList dto.MailingList) (dto.MailingList, error) {
	validator := validators.NewMailingListValidator()
	if err := validator.Validate(&newMailingList); err != nil {
		return newMailingList, err
	}

	mailingList := &dto.MailingList{
		Username:  newMailingList.Username,
		Email:     newMailingList.Email,
		CreatedAt: time.Now(),
	}

	repo := repositories.NewCsvMailingListRepository("mailing_list.csv")

	if err := repo.Save(mailingList); err != nil {
		return newMailingList, err
	}

	return *mailingList, nil
}
