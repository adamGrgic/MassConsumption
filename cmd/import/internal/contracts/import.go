package contracts

import "github.com/go-rod/rod"

type ImportRequest struct {
	SearchTerm string    `json: "searchTerm"`
	RodPage    *rod.Page `json: "rodPage"`
}
