package models

import (
	"time"

	"github.com/google/uuid"
)

type ImportProcess struct {
	Id          uuid.UUID  `json:"id"`
	ImportLogId uuid.UUID  `json:"importlogid"`
	Success     bool       `json:"success"`
	Error       string     `json:"error"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type ImportEntry struct {
	Id          uuid.UUID  `json:"id"`
	ImportLogId uuid.UUID  `json:"importlogid"`
	Success     bool       `json:"success"`
	Error       string     `json:"error"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
