package models

import "time"

type TaskList struct {
	Id          int        `json:"id"`
	Title       int        `json:"title"`
	Description int        `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
