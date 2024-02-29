package options

import (
	"fmt"
	"strings"
)

// Options contain filtering, pagination and sorting instructions provided via
// the querystring in bracketed object notation
type Options struct {
	ps IPaginationStrategy
	qs string

	Fields []string            `json:"fields,omitempty"`
	Filter map[string][]string `json:"filter,omitempty"`
	Page   map[string]int      `json:"page"`
	Sort   []string            `json:"sort,omitempty"`
}

// ContainsFilterField confirms whether or not the provided filter
// parameters include the requested field
func (o Options) ContainsFilterField(field string) bool {
	fields := []string{}
	for f := range o.Filter {
		fields = append(fields, f)
	}

	return contains(fields, field, false)
}

// ContainsSortField confirms whether or not the provided sort options
// contains the requested field
func (o Options) ContainsSortField(field string) bool {
	return contains(o.Sort, field, true)
}

// First returns a querystring for the first page
func (o Options) First() string {
	if len(o.Page) == 0 || o.ps == nil {
		return buildQuerystring(o.Filter, o.Fields, "", o.Sort)
	}

	// determine next page numbers based on pagination strategy
	po := o.ps.First(o.Page)
	qs := buildQuerystring(o.Filter, o.Fields, po, o.Sort)

	return qs
}

// Last returns a querystring for the last page
func (o Options) Last(total int) string {
	if len(o.Page) == 0 || o.ps == nil {
		return buildQuerystring(o.Filter, o.Fields, "", o.Sort)
	}

	// determine next page numbers based on pagination strategy
	po := o.ps.Last(o.Page, total)
	qs := buildQuerystring(o.Filter, o.Fields, po, o.Sort)

	return qs
}

// Next returns a querystring for the next page
func (o Options) Next() string {
	if len(o.Page) == 0 || o.ps == nil {
		return buildQuerystring(o.Filter, o.Fields, "", o.Sort)
	}

	// determine next page numbers based on pagination strategy
	po := o.ps.Next(o.Page)
	qs := buildQuerystring(o.Filter, o.Fields, po, o.Sort)

	return qs
}

// PaginationStrategy can be used to retrieve the current
// IPaginationStrategy that the Options struct will use for
// generating Prev, Next, First and Last querystring values
func (o Options) PaginationStrategy() IPaginationStrategy {
	return o.ps
}

// Prev returns a querystring for the previous page
func (o Options) Prev() string {
	if len(o.Page) == 0 || o.ps == nil {
		return buildQuerystring(o.Filter, o.Fields, "", o.Sort)
	}

	// determine previous page numbers based on pagination strategy
	po := o.ps.Prev(o.Page)
	qs := buildQuerystring(o.Filter, o.Fields, po, o.Sort)

	return qs
}

// String returns a querystring for the current page
func (o Options) String() string {
	return buildQuerystring(o.Filter, o.Fields, o.ps.Current(o.Page), o.Sort)
}

// SetPaginationStrategy can be used to specify custom pagination
// increments for Next, Prev, First and Last
func (o *Options) SetPaginationStrategy(ps IPaginationStrategy) {
	o.ps = ps
}

func buildQuerystring(filter map[string][]string, fields []string, page string, sort []string) string {
	b := strings.Builder{}
	ra := false

	// filters
	for field, filter := range filter {
		if ra {
			fmt.Fprint(&b, "&")
		}

		// & is required on subsequent iterations
		ra = true

		// write the filter for the field to the builder
		fmt.Fprintf(&b, "filter[%s]=", field)
		for i, value := range filter {
			// add a comma if multiple values are specified
			if i > 0 {
				fmt.Fprint(&b, ",")
			}
			fmt.Fprint(&b, value)
		}
	}

	// field projections
	if len(fields) > 0 {
		if ra {
			fmt.Fprint(&b, "&")
		}
		fmt.Fprintf(&b, "fields=")
		for i, field := range fields {
			// add a comma if multiple fields are specified
			if i > 0 {
				fmt.Fprint(&b, ",")
			}

			fmt.Fprint(&b, field)
		}
	}

	// pagination
	if page != "" {
		if ra {
			fmt.Fprint(&b, "&")
		}

		// & is required on subsequent iterations
		ra = true

		fmt.Fprint(&b, page)
	}

	// sorting
	if len(sort) > 0 {
		if ra {
			fmt.Fprint(&b, "&")
		}
		fmt.Fprintf(&b, "sort=")
		for i, field := range sort {
			// add a comma if multiple fields are specified
			if i > 0 {
				fmt.Fprint(&b, ",")
			}

			fmt.Fprint(&b, field)
		}
	}

	return b.String()
}

func contains(list []string, value string, stripPrefix bool) bool {
	if len(list) == 0 {
		return false
	}

	if value == "" {
		return false
	}

	for _, search := range list {
		// check to see if prefix should be stripped
		if stripPrefix {
			if len(search) > 2 && (search[0:2] == "<=" || search[0:2] == ">=" || search[0:2] == "!=") {
				search = search[2:]
			}

			if len(search) > 1 && (search[0:1] == "<" || search[0:1] == ">" || search[0:1] == "-" || search[0:1] == "+" || search[0:1] == "!") {
				search = search[1:]
			}
		}

		if search == value {
			return true
		}
	}

	return false
}
