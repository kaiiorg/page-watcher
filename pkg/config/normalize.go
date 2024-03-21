package config

import "regexp"

type Normalize struct {
	Regex string `hcl:"regex"`
	To    string `hcl:"to"`

	r *regexp.Regexp
}
