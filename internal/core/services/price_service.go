package services

import (
	"encoding/json"
	"fmt"
	"os"
	"web-scraper/internal/core/models"
)

type PriceService interface {
	GetPrices() ([]models.PriceRecord, error)
}

type PriceServiceImpl struct {
}

func NewPriceService() PriceService {
	return &PriceServiceImpl{}
}

func (s *PriceServiceImpl) GetPrices() ([]models.PriceRecord, error) {

	var p []models.PriceRecord
	importFile := os.Getenv("PRICES_IMPORT")

	_, err := os.Stat(importFile)
	if err != nil {
		fmt.Println("no import file available")
		return p, err
	}

	f, err := os.Open(importFile)
	if err != nil {
		fmt.Println("Could not open file: ", err)
		return p, err
	}
	defer f.Close()

	bytes, err := os.ReadFile(importFile)
	if err != nil {
		fmt.Println("could not read file ", err)
		return p, err
	}

	json.Unmarshal(bytes, &p)

	return p, nil
}
