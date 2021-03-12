// Package queryoptions provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package queryoptions

import (
	"regexp"
)

var (
	filterRE = regexp.MustCompile(`(?i)(?P<fc>mandatory|optional)\]\[(?P<e>\w+)\]\[(?P<f>\w+)\]\=`)
)

type Options struct {
	Filter FilterContainer `json:"filter"`
	Page   Page            `json:"page"`
	Sort   []string        `json:"sort"`
}

func FromQuerystring(qs string) (*Options, error) {
	options := &Options{}

	if qs == "" {
		return options, nil
	}

	// look for filters
	if err := parseFilters(qs, options); err != nil {
		return nil, err
	}

	return options, nil
}

func parseFilters(qs string, options *Options) error {
	// check for filters
	if !filterRE.MatchString(qs) {
		return nil
	}

	pm := make(map[string]string)
	m := filterRE.FindStringSubmatch(qs)
	for i, n := range filterRE.SubexpNames() {
		if i > 0 && i <= len(m) {
			pm[n] = m[i]
		}
	}

	return nil
}
