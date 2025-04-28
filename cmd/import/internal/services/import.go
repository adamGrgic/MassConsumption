package services

import (
	"fmt"
	"time"
	"web-scraper/cmd/import/internal/contracts"
	"web-scraper/internal/core/models"

	"github.com/rs/zerolog/log"
)

// maybe change to Collector
type Transformer struct {
	FailChannel    chan<- string             `json:"failChannel"`
	SuccessChannel chan<- string             `json:"successChannel"`
	PricesChannel  chan<- models.PriceRecord `json:"pricesChannel"`
}

func NewTransformer(success chan<- string, fail chan<- string, prices chan<- models.PriceRecord) *Transformer {
	return &Transformer{SuccessChannel: success, FailChannel: fail, PricesChannel: prices}
}

func (r *Transformer) Process(ii contracts.ImportRequest) error {
	elements := ii.RodPage.MustElements(`div[data-test="product-grid"] div[data-test="product-body"]`)

	if len(elements) == 0 {
		log.Warn().Str("search term", ii.SearchTerm).Msg("Could not find element for given search term: ")
		select {
		case r.FailChannel <- ii.SearchTerm:
		case <-time.After(5 * time.Second):
			log.Warn().Str("search term", ii.SearchTerm).Msg("Timeout trying to send to failChan")
		}
		return fmt.Errorf("no product grid elements found for search term: %s", ii.SearchTerm)
	} else {
		fmt.Println(elements.First().Element("div"))
	}

	// _, err := elements.First().Element(`div[data-test="product-title"]`)
	// if err != nil {
	// 	log.Warn().Str("search term", ii.SearchTerm).Msg("Could not find element for given search term: ")
	// 	select {
	// 	case r.FailChannel <- ii.SearchTerm:
	// 	case <-time.After(5 * time.Second):
	// 		log.Warn().Str("search term", ii.SearchTerm).Msg("Timeout trying to send to failChan")
	// 	}
	// 	return err
	// }

	for _, el := range elements {
		titleEl, err := el.Element(`div[data-test="product-title"]`)
		if err != nil {
			log.Err(err).Msg("Couldn't get title element")
			continue
		}
		priceEl, err := el.Element(`div[data-test="current-price"]`) // <-- small fix: you used wrong selector
		if err != nil {
			log.Err(err).Msg("Couldn't get price element")
			continue
		}

		// Now extract text safely
		titleText, err := titleEl.Text()
		if err != nil {
			log.Err(err).Msg("Couldn't get title text")
			continue
		}

		priceText, err := priceEl.Text()
		if err != nil {
			log.Err(err).Msg("Couldn't get price text")
			continue
		}

		fmt.Println("Title:", titleText)
		fmt.Println("Price:", priceText)

		r.PricesChannel <- models.PriceRecord{Title: titleText, Price: priceText, Category: ii.SearchTerm}

	}

	select {
	case r.SuccessChannel <- ii.SearchTerm:
	case <-time.After(5 * time.Second):
		log.Warn().Str("search term", ii.SearchTerm).Msg("Timeout trying to send to prices channel")
	}

	return nil
}
