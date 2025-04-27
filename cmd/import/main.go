package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"web-scraper/internal/core/models"

	"github.com/rs/zerolog/log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var searchTerms = []string{"peanut butter", "banana", "nutella", "ground beef", "chicken", "eggs"}

func main() {
	start := time.Now()

	l := launcher.New().NoSandbox(true).Headless(true).Devtools(false)
	url := l.MustLaunch()

	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	var wg sync.WaitGroup
	resultsChan := make(chan models.PriceRecord)
	failChan := make(chan string, 1000)
	// failChan := make(chan string)

	for _, st := range searchTerms {
		wg.Add(1)
		go func(searchTerm string) {
			defer wg.Done()

			incognito := browser.MustIncognito()
			defer incognito.MustClose()

			link := fmt.Sprintf("https://www.target.com/s?searchTerm=%s", strings.ReplaceAll(searchTerm, " ", "+"))
			page := incognito.MustPage(link)

			page.MustWaitLoad()
			page = page.Timeout(15 * time.Second)
			page.MustElement(`div[data-test="product-grid"]`)
			fmt.Println("after must element")
			page.MustWaitIdle()

			elements := page.MustElements(`div[data-test="product-grid"] > *`)

			for _, el := range elements {
				titleEl, err := el.Element(`div[data-test="product-title"]`)
				if err != nil {
					log.Err(err).Str("search term", st).Msg("Could not find title element for given search term.")
					select {
					case failChan <- searchTerm:
					case <-time.After(5 * time.Second):
						log.Warn().Str("search term", searchTerm).Msg("Timeout trying to send to failChan")
					}
					continue
				}

				priceEl, err := el.Element(`span[data-test="current-price"]`)
				if err != nil {
					log.Err(err).Str("search term", st).Msg("Could not find title element for given search term.")
					select {
					case failChan <- searchTerm:
					case <-time.After(5 * time.Second):
						log.Warn().Str("search term", searchTerm).Msg("Timeout trying to send to failChan")
					}
					continue
				}

				title := titleEl.MustText()
				price := priceEl.MustText()

				resultsChan <- models.PriceRecord{Title: title, Price: price, Category: searchTerm}
			}
		}(st)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(failChan)
	}()

	var priceResult []models.PriceRecord
	for r := range resultsChan {
		priceResult = append(priceResult, r)
	}

	var failedSearches []string
	for r := range failChan {
		failedSearches = append(failedSearches, r)
	}

	if len(priceResult) == 0 {
		log.Warn().Msg("No price records collected, skipping file write")
		return
	}

	err := os.MkdirAll("./tmp", os.ModePerm)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create tmp directory")
	}

	finalData, err := json.MarshalIndent(priceResult, "", "  ")
	if err != nil {
		log.Panic().Err(err).Msg("Could not marshal price data")
	}

	err = os.WriteFile("./tmp/prices.json", finalData, 0644)
	if err != nil {
		log.Err(err).Msg("Failed to write prices file")
		return
	}

	fmt.Println("foo")

	log.Info().
		Int("record_count", len(priceResult)).
		Dur("duration", time.Since(start)).
		Msg("Scraping completed successfully")

}
