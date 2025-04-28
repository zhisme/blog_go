package repositories

import (
	"encoding/csv"
	"os"
	"time"
	"backend-go/internal/dto"
)

type CsvMailingListRepository struct {
	filepath string
}

func NewCsvMailingListRepository(filepath string) *CsvMailingListRepository {
	return &CsvMailingListRepository{
		filepath: filepath,
	}
}

func (r *CsvMailingListRepository) Save(mailingList *dto.MailingList) error {
	file, err := os.OpenFile(r.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		headers := []string{"Username", "Email", "CreatedAt"}
		if err := writer.Write(headers); err != nil {
			return err
		}
	}

	createdAt := mailingList.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	record := []string{
		mailingList.Username,
		mailingList.Email,
		createdAt.Format(time.RFC3339),
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	return writer.Error()
}
