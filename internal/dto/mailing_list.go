package dto

import "time"

type MailingList struct {
	CreatedAt time.Time `json:"createdAt"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}
