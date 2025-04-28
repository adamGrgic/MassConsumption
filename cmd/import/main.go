package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"time"
	"web-scraper/cmd/import/internal/contracts"
	"web-scraper/cmd/import/internal/services"
	"web-scraper/internal/core/models"
	"web-scraper/internal/logging"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

var searchTerms = []string{"peanut butter", "banana", "nutella", "ground beef", "chicken", "eggs"}

func main() {
	start := time.Now()

	godotenv.Load()

	logging.ConfigureLogging()

	l := launcher.New().NoSandbox(true).Headless(true).Devtools(false)
	url := l.MustLaunch()

	browser := rod.New().SlowMotion(time.Duration(rand.IntN(100)) * time.Millisecond).ControlURL(url).MustConnect()
	defer browser.MustClose()

	var wg sync.WaitGroup
	pc := make(chan models.PriceRecord)
	fc := make(chan string, 100)
	sc := make(chan string, 100)

	// retries := 5

	t := services.NewTransformer(sc, fc, pc)

	for _, st := range searchTerms {
		wg.Add(1)
		go func(searchTerm string) {
			defer wg.Done()

			incognito := browser.MustIncognito()
			defer incognito.MustClose()

			link := fmt.Sprintf("https://www.target.com/s?searchTerm=%s", strings.ReplaceAll(searchTerm, " ", "+"))
			page := incognito.MustPage(link)
			page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
			})

			_, err := page.HTML()
			if err != nil {
				log.Err(err).Msg("Could not get page HTML")
			}

			page.MustWaitLoad()
			page = page.Timeout(15 * time.Second)
			page.MustElement(`div[data-test="product-grid"]`)
			page.MustWaitIdle()

			err = t.Process(contracts.ImportRequest{SearchTerm: searchTerm, RodPage: page})

		}(st)
	}

	go func() {
		wg.Wait()
		close(sc)
		close(fc)
		close(pc)
	}()

	var priceResult []models.PriceRecord
	for r := range pc {
		priceResult = append(priceResult, r)
	}

	var failedSearches []string
	for f := range fc {
		failedSearches = append(failedSearches, f)
	}

	var successSearches []string
	for s := range sc {
		successSearches = append(successSearches, s)
	}

	if len(priceResult) == 0 {
		log.Warn().Msg("No price records collected, skipping file write")
		return
	}

	err := os.MkdirAll("./tmp", os.ModePerm)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create tmp directory")
	}

	fr, err := json.MarshalIndent(priceResult, "", "  ")
	if err != nil {
		log.Panic().Err(err).Msg("Could not marshal price data")
	}

	err = os.WriteFile("./tmp/prices.json", fr, 0644)
	if err != nil {
		log.Err(err).Msg("Failed to write prices file")
		return
	}

	log.Info().
		Int("imported prices", len(priceResult)).
		Int("success count", len(successSearches)).
		Int("failed count", len(failedSearches)).
		Dur("duration", time.Since(start)).
		Msg("Scraping completed successfully")

}
