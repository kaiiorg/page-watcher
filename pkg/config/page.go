package config

import "time"

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
}

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
