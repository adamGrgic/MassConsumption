package models

type PriceRecord struct {
	// ID    uuid.UUID `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	// CategoryId uuid.UUID `json:"categoryid"`
	Price string `json:"price"`
	// CurrencyId string    `json:"currencyid"`
	// ProcessId uuid.UUID `json:"processid"`
	// CreatedAt time.Time `json:"created_at"`
}
