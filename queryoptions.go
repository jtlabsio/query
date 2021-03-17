// Package queryoptions provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package queryoptions

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	bracketRE = regexp.MustCompile(`(?P<typ>filter|sort|page)\[(.+?)\](\=?)`)
	sortRE    = regexp.MustCompile(`sort=(?P<field>.+?)(\&|\z)`)
	valueRE   = regexp.MustCompile(`\=(.+?)(\&|\z)`)
)

// Options contain filtering, pagination and sorting instructions provided via
// the querystring in bracketed object notation
type Options struct {
	Filter map[string]string `json:"filter"`
	Page   map[string]int    `json:"page"`
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
	options.Sort = parseSort(qs)

	return options, nil
}

func parseBracketParams(qs string) (Options, error) {
	o := Options{
		Filter: map[string]string{},
		Page:   map[string]int{},
		Sort:   []string{},
	}
	terms := bracketRE.FindAllStringSubmatch(qs, -1)
	values := valueRE.FindAllStringSubmatch(qs, -1)

	if len(terms) > 0 && len(terms) != len(values) {
		// multiple nested bracket params... not sure how to parse
		return o, errors.New("unable to parse: an object hierarchy has been provided")
	}

	for i, term := range terms {
		switch strings.ToLower(term[1]) {
		case "filter":
			if o.Filter == nil {
				o.Filter = map[string]string{}
			}
			o.Filter[term[2]] = values[i][1]
		case "page":
			if o.Page == nil {
				o.Page = map[string]int{}
			}

			v, err := strconv.ParseInt(values[i][1], 0, 64)
			if err != nil {
				return o, err
			}

			o.Page[term[2]] = int(v)
		}
	}

	return o, nil
}

func parseSort(qs string) []string {
	sort := []string{}
	fields := sortRE.FindAllStringSubmatch(qs, -1)

	for _, field := range fields {
		sort = append(sort, field[0])
	}

	return sort
}
