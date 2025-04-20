package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StatusId    int        `json:"status"`
	ListId      int        `json:"listid"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
