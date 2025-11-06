package repositories

import (
	"backend-go/internal/dto"
	"encoding/csv"
	"log"
	"os"
	"time"
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
	exists, err := r.emailExists(mailingList.Email)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Email already subscribed: %s", mailingList.Email)
		return nil
	}

	file, err := os.OpenFile(r.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

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

func (r *CsvMailingListRepository) emailExists(email string) (bool, error) {
	file, err := os.Open(r.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	// Skip header: start from index 1
	for i, record := range records {
		if i == 0 {
			continue
		}
		if record[1] == email {
			return true, nil
		}
	}

	return false, nil
}
