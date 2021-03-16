// Package queryoptions provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package queryoptions

import (
	"errors"
	"regexp"
	"strings"
)

var (
	bracketRE = regexp.MustCompile(`(?P<typ>filter|sort|page)\[(.+?)\](\=?)`)
	valueRE   = regexp.MustCompile(`\=(.+?)(\&|\z)`)
)

// Options contain filtering, pagination and sorting instructions provided via
// the querystring in bracketed object notation
type Options struct {
	Filter map[string]string `json:"filter"`
	Page   map[string]string `json:"page"`
	Sort   []string          `json:"sort"`
}

// FromQuerystring parses an Options object from the provided querystring
func FromQuerystring(qs string) (Options, error) {
	if qs == "" {
		return Options{}, nil
	}

	// apply filter and page
	options, err := parseBracketParams(qs)
	if err != nil {
		return options, err
	}

	// apply sort

	return options, nil
}

func parseBracketParams(qs string) (Options, error) {
	o := Options{}
	terms := bracketRE.FindAllStringSubmatch(qs, -1)
	values := valueRE.FindAllStringSubmatch(qs, -1)

	if len(terms) != len(values) {
		// multiple nested bracket params... not sure how to parse
		return o, errors.New("unable to parse: an object hierarchy has been provided")
	}

	for i, term := range terms {
		switch strings.ToLower(term[1]) {
		case "filter":
			if o.Filter == nil {
				o.Filter = make(map[string]string)
			}
			o.Filter[term[2]] = values[i][1]
		case "page":
			if o.Page == nil {
				o.Page = make(map[string]string)
			}
			o.Page[term[2]] = values[i][1]
		}
	}

	return o, nil
}
