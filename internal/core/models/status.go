package models

type Status struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Label string `json: "label"`
}
