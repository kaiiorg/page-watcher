package config

import (
	"github.com/rs/zerolog/log"
	"regexp"
	"time"
)

const (
	DefaultPageEvery = time.Minute
	MinimumPageEvery = 5 * time.Second
)

type Page struct {
	// Name is a friendly name for this page. Used in log messages
	Name string `hcl:"name"`
	// Url is the URL of the web page to be watched
	Url string `hcl:"url"`
	// Every is a string duration for how often to check the Url. Rounded to nearest second, min 5s, default 1m
	Every string `hcl:"every"`

	// Find is a string slice for the arguments to soup's Root.Find(). See https://pkg.go.dev/github.com/anaskhan96/soup#Root.Find
	Find []string `hcl:"find"`

	// Normalize is a map of regular expression strings to a string that will replace anything that matches
	// Hint: use https://regex101.com/ or something similar, select Golang flavor
	Normalize map[string]string `hcl:"normalize,optional"`
	// normalizeRegex is the compiled regexps defined in Normalize
	normalizeRegex map[*regexp.Regexp]string

	// Debug will write the page to files to help with configuration
	Debug bool `hcl:"debug,optional"`
}

// EveryDuration parses the configured duration string and returns a default value if it is invalid
func (p *Page) EveryDuration() time.Duration {
	d, err := time.ParseDuration(p.Every)
	if err == nil {
		d = d.Truncate(time.Second)
		if d < MinimumPageEvery {
			return MinimumPageEvery
		}
		return d
	}
	return DefaultPageEvery
}

// ValidateNormalize attempts to compile the configured regex map to a valid regex
func (p *Page) ValidateNormalize() error {
	rMap := map[*regexp.Regexp]string{}

	for expr, to := range p.Normalize {
		r, err := regexp.Compile(expr)
		if err != nil {
			return err
		}
		rMap[r] = to
	}

	p.normalizeRegex = rMap

	return nil
}

// NormalizeString will run the regex defined in Page.Normalize
func (p *Page) NormalizeString(s string) string {
	for r, to := range p.normalizeRegex {
		log.Debug().Str("configKey", r.String()).Msg("Running configured normalize regex")
		s = r.ReplaceAllString(s, to)
	}
	return s
}
