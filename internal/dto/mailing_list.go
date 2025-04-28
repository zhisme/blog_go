package dto

import "time"

type MailingList struct {
  Username   string    `json:"username"`
  Email      string    `json:"email"`
  CreatedAt  time.Time `json:"createdAt"`
}
