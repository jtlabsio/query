// options provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package options

import (
	"errors"
	"fmt"
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

	options := Options{
		qs: uqs,
	}

	// parse fields
	options.Fields = parseFields(&uqs)

	// parse sort
	options.Sort = parseSort(&uqs)

	// parse filter and page
	if err := parseBracketParams(uqs, &options); err != nil {
		return options, err
	}

	// attempt to infer pagination strategy
	if _, ok := options.Page["limit"]; ok {
		options.SetPaginationStrategy(&OffsetStrategy{})
	}

	if _, ok := options.Page["size"]; ok {
		options.SetPaginationStrategy(&PageSizeStrategy{})
	}

	return options, nil
}

func extract(qs *string, re regexp.Regexp) string {
	coords := re.FindStringIndex(*qs)

	var r string
	if coords != nil {
		if coords[0] == 0 {
			// at beginning of string
			r = (*qs)[coords[0]:coords[1]]
			*qs = (*qs)[coords[1]:]
			return r
		}

		if coords[1] == len(*qs) {
			// at end of string
			r = (*qs)[coords[0]:coords[1]]
			*qs = (*qs)[0:coords[0]]
			return r
		}

		// in middle of string
		r = (*qs)[coords[0]:coords[1]]
		*qs = fmt.Sprintf("%s%s", (*qs)[0:coords[0]], (*qs)[coords[1]:])
		return r
	}

	return r
}

func parseBracketParams(qs string, o *Options) error {
	o.Filter = map[string][]string{}
	o.Page = map[string]int{}

	terms := bracketRE.FindAllStringSubmatch(qs, -1)
	values := valueRE.FindAllStringSubmatch(qs, -1)

	if len(terms) > 0 && len(terms) > len(values) {
		// multiple nested bracket params... not sure how to parse
		return errors.New("unable to parse: an object hierarchy has been provided")
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
				return err
			}

			o.Page[term[2]] = int(v)
		}
	}

	return nil
}

func parseFields(qs *string) []string {
	fields := []string{}

	fieldNames := fieldsRE.FindAllStringSubmatch(extract(qs, *fieldsRE), -1)
	for fieldNames != nil {
		for _, field := range fieldNames {
			// check if sort value is an array
			if commaRE.MatchString(field[1]) {
				fa := commaRE.Split(field[1], -1)
				fields = append(fields, fa...)
				continue
			}

			fields = append(fields, field[1])
		}

		// look for more fields= occurrences
		fieldNames = fieldsRE.FindAllStringSubmatch(extract(qs, *fieldsRE), -1)
	}

	return fields
}

func parseSort(qs *string) []string {
	sort := []string{}

	fieldNames := sortRE.FindAllStringSubmatch(extract(qs, *sortRE), -1)
	for fieldNames != nil {
		for _, field := range fieldNames {
			// check if sort value is an array
			if commaRE.MatchString(field[1]) {
				fa := commaRE.Split(field[1], -1)
				sort = append(sort, fa...)
				continue
			}

			sort = append(sort, field[1])
		}

		// look for more sort= occurrences
		fieldNames = sortRE.FindAllStringSubmatch(extract(qs, *sortRE), -1)
	}

	return sort
}
