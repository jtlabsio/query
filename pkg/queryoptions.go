// Package queryoptions provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package queryoptions

import (
	"regexp"
)

var (
	filterRE     = regexp.MustCompile(`&?filter\[`)
	filterTypeRE = regexp.MustCompile(`(?i)(?P<type>mandatory|optional)\]\[(?P<filter>\w+)\]\[(?P<name>\w+)\]\=(?P<val>[\w\d\D]+)`)
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

	filters := filterRE.Split(qs, -1)

	for _, f := range filters {
		if f == "" {
			continue
		}

		var fm *map[string]string
		m := filterTypeRE.FindStringSubmatch(f)

		// identify the correct map (mandatory)
		if m[1] == "mandatory" {
			switch m[2] {
			case "beginsWith":
				fm = &options.Filter.Mandatory.BeginsWith
			case "contains":
				fm = &options.Filter.Mandatory.Contains
			case "endsWith":
				fm = &options.Filter.Mandatory.EndsWith
			case "exact":
				fm = &options.Filter.Mandatory.Exact
			case "greaterThan":
				fm = &options.Filter.Mandatory.GreaterThan
			case "greaterThanEqual":
				fm = &options.Filter.Mandatory.GreaterThanEqual
			case "gt":
				fm = &options.Filter.Mandatory.LessThanEqual
			case "gte":
				fm = &options.Filter.Mandatory.LessThanEqual
			case "lessThan":
				fm = &options.Filter.Mandatory.LessThan
			case "lessThanEqual":
				fm = &options.Filter.Mandatory.LessThanEqual
			case "lt":
				fm = &options.Filter.Mandatory.LessThanEqual
			case "lte":
				fm = &options.Filter.Mandatory.LessThanEqual
			case "ne":
				fm = &options.Filter.Mandatory.NotEqual
			case "notEqual":
				fm = &options.Filter.Mandatory.NotEqual
			}
		}

		// identify the correct map (optional)
		if m[1] == "optional" {
			switch m[2] {
			case "beginsWith":
				fm = &options.Filter.Optional.BeginsWith
			case "contains":
				fm = &options.Filter.Optional.Contains
			case "endsWith":
				fm = &options.Filter.Optional.EndsWith
			case "exact":
				fm = &options.Filter.Optional.Exact
			case "greaterThan":
				fm = &options.Filter.Optional.GreaterThan
			case "greaterThanEqual":
				fm = &options.Filter.Optional.GreaterThanEqual
			case "gt":
				fm = &options.Filter.Optional.LessThanEqual
			case "gte":
				fm = &options.Filter.Optional.LessThanEqual
			case "lessThan":
				fm = &options.Filter.Optional.LessThan
			case "lessThanEqual":
				fm = &options.Filter.Optional.LessThanEqual
			case "lt":
				fm = &options.Filter.Optional.LessThanEqual
			case "lte":
				fm = &options.Filter.Optional.LessThanEqual
			case "ne":
				fm = &options.Filter.Optional.NotEqual
			case "notEqual":
				fm = &options.Filter.Optional.NotEqual
			}
		}

		// if map doesn't exist, create it...
		if (*fm) == nil {
			(*fm) = map[string]string{m[3]: m[4]}
			continue
		}

		// assign the provided field name and value to the filter map
		(*fm)[m[3]] = m[4]
	}

	return nil
}
