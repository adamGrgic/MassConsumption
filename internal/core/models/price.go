package models

import (
	"time"

	"github.com/google/uuid"
)

type PriceRecord struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	CategoryId uuid.UUID `json:"categoryid"`
	Price      int64     `json:"price"`
	CurrencyId string    `json:"currencyid"`
	ProcessId  uuid.UUID `json:"processid"`
	CreatedAt  time.Time `json:"created_at"`
}
