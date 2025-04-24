package main

// comment new foo
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

type PriceResult struct {
	Category string `json:"category"`
	Price    string `json:"price"`
	Title    string `json:"title"`
}

var searchTerms = []string{"peanut butter", "banana", "nutella", "ground beef", "chicken", "eggs", ""}

func main() {
	var priceResult []PriceResult

	// Headless runs the browser on foreground, you can also use flag "-rod=show"
	// Devtools opens the tab in each new tab opened automatically
	l := launcher.New().
		NoSandbox(true).
		Headless(true).
		Devtools(true)

	defer l.Cleanup()

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(2 * time.Second).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()

	for _, st := range searchTerms {
		link := fmt.Sprintf("https://www.target.com/s?searchTerm=%s", strings.ReplaceAll(st, " ", "+"))
		page := browser.MustPage(link)

		time.Sleep(10 * time.Second)

		page.MustWaitIdle()
		page.MustElement(`div[data-test="product-grid"]`)

		elements := page.MustElements(`div[data-test="product-grid"] > *`)

		fmt.Println(elements)

		if elements != nil {
			for i, el := range elements {
				titleEl, titleErr := el.Element(`div[data-test="product-title"]`)
				if titleErr != nil {
					fmt.Println("skipping invalid search result")
					continue
				}

				priceEl, priceErr := el.Element(`span[data-test="current-price"]`)
				if priceErr != nil {
					fmt.Println("skipping invalid search result")
					continue
				}

				if titleEl == nil || priceEl == nil {
					continue // skip if it's not a valid product card
				}

				title := titleEl.MustText()
				price := priceEl.MustText()

				priceResult = append(priceResult, PriceResult{Title: title, Price: price, Category: st})

				fmt.Printf("Item #%d: %s - %s\n", i+1, title, price)
			}
		}

		fmt.Println("Final Price Result: ", priceResult)

		finalData, err := json.MarshalIndent(priceResult, "", "  ")
		if err != nil {
			log.Panic("Could not marshal price data: ", err)
		}

		os.WriteFile("./tmp/prices.json", finalData, 0644)

	}

	utils.Pause() // pause goroutine
}
