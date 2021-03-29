// Package queryoptions provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package queryoptions

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	bracketRE = regexp.MustCompile(`(?P<typ>filter|sort|page)\[(.+?)\](\=?)`)
	commaRE   = regexp.MustCompile(`\s?\,\s?`)
	fieldsRE  = regexp.MustCompile(`fields=(?P<field>.+?)(\&|\z)`)
	sortRE    = regexp.MustCompile(`sort=(?P<field>.+?)(\&|\z)`)
	valueRE   = regexp.MustCompile(`\=(.+?)(\&|\z)`)
)

// FromQuerystring parses an Options object from the provided querystring
func FromQuerystring(qs string) (Options, error) {
	if qs == "" {
		return Options{}, nil
	}

	uqs, err := url.QueryUnescape(qs)
	if err != nil {
		return Options{}, err
	}

	// apply filter and page
	options, err := parseBracketParams(uqs)
	if err != nil {
		return options, err
	}

	// apply fields
	options.Fields = parseFields(qs)

	// apply sort
	options.Sort = parseSort(qs)

	// attempt to infer pagination strategy
	if _, ok := options.Page["limit"]; ok {
		options.SetPaginationStrategy(&OffsetStrategy{})
	}

	if _, ok := options.Page["size"]; ok {
		options.SetPaginationStrategy(&PageSizeStrategy{})
	}

	return options, nil
}

func parseBracketParams(qs string) (Options, error) {
	o := Options{
		qs:     qs,
		Filter: map[string][]string{},
		Page:   map[string]int{},
		Sort:   []string{},
	}
	terms := bracketRE.FindAllStringSubmatch(qs, -1)
	values := valueRE.FindAllStringSubmatch(qs, -1)

	if len(terms) > 0 && len(terms) > len(values) {
		// multiple nested bracket params... not sure how to parse
		return o, errors.New("unable to parse: an object hierarchy has been provided")
	}

	for i, term := range terms {
		switch strings.ToLower(term[1]) {
		case "filter":
			if o.Filter == nil {
				o.Filter = map[string][]string{}
			}

			// check for array
			if commaRE.MatchString(values[i][1]) {
				o.Filter[term[2]] = commaRE.Split(values[i][1], -1)
				continue
			}

			o.Filter[term[2]] = []string{values[i][1]}
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

func parseFields(qs string) []string {
	fields := []string{}
	fieldNames := fieldsRE.FindAllStringSubmatch(qs, -1)

	for _, field := range fieldNames {
		// check if sort value is an array
		if commaRE.MatchString(field[1]) {
			fa := commaRE.Split(field[1], -1)
			fields = append(fields, fa...)
			continue
		}

		fields = append(fields, field[1])
	}

	return fields
}

func parseSort(qs string) []string {
	sort := []string{}
	fields := sortRE.FindAllStringSubmatch(qs, -1)

	for _, field := range fields {
		// check if sort value is an array
		if commaRE.MatchString(field[1]) {
			fa := commaRE.Split(field[1], -1)
			sort = append(sort, fa...)
			continue
		}

		sort = append(sort, field[1])
	}

	return sort
}
