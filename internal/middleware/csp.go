package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Nonces struct {
	Htmx            string
	ResponseTargets string
	CSS             string
	HtmxCSSHash     string
	JS              string
}

const nonceKey = "nonces"

func generateNonce() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate nonce: %v", err)
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		nonces := Nonces{
			Htmx:            generateNonce(),
			ResponseTargets: generateNonce(),
			CSS:             generateNonce(),
			JS:              generateNonce(),
			HtmxCSSHash:     "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg=",
		}

		// Store nonces in Gin context
		c.Set(nonceKey, nonces)

		// Construct CSP header
		cspHeader := fmt.Sprintf(
			"default-src 'self'; "+
				"script-src 'self' 'nonce-%s' 'nonce-%s'; "+
				"style-src 'self' 'nonce-%s' '%s'; "+
				"object-src 'none'; "+
				"base-uri 'none';",
			nonces.Htmx,
			nonces.JS,
			nonces.CSS,
			nonces.HtmxCSSHash,
		)
		c.Header("Content-Security-Policy", cspHeader)
		c.Header("Content-Security-Policy-Report-Only", cspHeader)

		c.Next()
	}
}

// Helper functions to retrieve nonces from Gin context

func GetNonces(c *gin.Context) Nonces {
	val, exists := c.Get(nonceKey)
	if !exists {
		log.Fatal("Nonces missing in context - did you apply CSPMiddleware?")
	}

	nonces, ok := val.(Nonces)
	if !ok {
		log.Fatal("Invalid nonce type stored in context")
	}

	return nonces
}

func GetHtmxNonce(c *gin.Context) string {
	return GetNonces(c).Htmx
}

func GetResponseTargetsNonce(c *gin.Context) string {
	return GetNonces(c).ResponseTargets
}

func GetCssNonce(c *gin.Context) string {
	return GetNonces(c).CSS
}

func GetJsNonce(c *gin.Context) string {
	return GetNonces(c).JS
}
