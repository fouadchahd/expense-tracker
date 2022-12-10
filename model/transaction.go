package model

import (
	"errors"
	"time"
)

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	SubmittedAt time.Time `json:"submitted_at"`
	UserId      uint      `json:"user_id"`
	CategoryId  uint      `json:"category_id"`
	GroupID     *uint     `json:"group_id"`
}

func (t *Transaction) IsValid() (bool, error) {
	if t.SubmittedAt.IsZero() {
		now := time.Now().UTC()
		t.SubmittedAt = now
	}

	if t.Amount == 0 {
		return false, errors.New("invalid amount")
	}
	if t.UserId == 0 {
		return false, errors.New("missing user")
	}
	if !(t.CategoryId > 0) {
		return false, errors.New("missing category")
	}
	return true, nil
}
