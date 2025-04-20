package media

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func GetCSSFile(file string) string {

	manifestPath := os.Getenv("MEDIA_MANIFEST_PATH")

	absolutePath, err := filepath.Abs(manifestPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("manifestPath", manifestPath).
			Str("absolutePath", absolutePath).
			Msg("Unable to retrieve media manifest file")
	}

	f, err := os.Open(absolutePath)
	if err != nil {
		log.Fatal().Msg("This is a fatal log")
	}
	defer f.Close()

	var manifest map[string]string
	if err := json.NewDecoder(f).Decode(&manifest); err != nil {
		log.Fatal().Msg("This is a fatal log")
	}

	key := "css:" + file
	val, ok := manifest[key]
	if !ok {
		fmt.Fprintf(os.Stderr, "❌ Key '%s' not found in manifest\n", key)
		os.Exit(1)
	}

	return val
}
