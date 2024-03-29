package normalizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaiiorg/page-watcher/pkg/config"

	"github.com/anaskhan96/soup"
	"github.com/rs/zerolog/log"
)

type PageNormalizer struct{}

func New() *PageNormalizer {
	return &PageNormalizer{}
}

func (pn *PageNormalizer) Get(page *config.Page) (string, string, error) {
	resp, err := soup.Get(page.Url)
	if err != nil {
		return "", "", err
	}

	doc := soup.HTMLParse(resp)
	content := doc.Find(page.Find...)
	if content.Error != nil {
		return "", "", content.Error
	}

	// Run user defined normalization regexp
	normalized := page.NormalizeString(content.FullText())

	if page.Debug {
		pn.debug(page.Name, resp, normalized)
	}

	return content.FullText(), normalized, nil
}

func (pn *PageNormalizer) debug(name, raw, normalized string) {
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "'", "")

	rawFileName := fmt.Sprintf("raw.%s.html", name)
	normalizedFileName := fmt.Sprintf("normalized.%s.txt", name)

	err := os.WriteFile(
		filepath.Join(".", "debug", rawFileName),
		[]byte(raw),
		0777,
	)
	if err != nil {
		log.Warn().Str("file", rawFileName).Msg("Failed to write raw contents to file")
	}

	err = os.WriteFile(
		filepath.Join(".", "debug", normalizedFileName),
		[]byte(normalized),
		0777,
	)
	if err != nil {
		log.Warn().Str("file", normalizedFileName).Msg("Failed to write normalized contents to file")
	}
}
