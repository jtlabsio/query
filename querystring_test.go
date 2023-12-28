// options provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package options

import (
	"reflect"
	"testing"
)

func TestFromQuerystring(t *testing.T) {
	offsetPS := &OffsetStrategy{}
	type args struct {
		ps IPaginationStrategy
		qs string
	}
	tests := []struct {
		name    string
		args    args
		want    Options
		wantErr bool
	}{
		{
			"empty querystring",
			args{qs: ""},
			Options{},
			false,
		},
		{
			"multiple filters not repeated",
			args{qs: "filter[fieldA]=value1&filter[fieldB]=value2"},
			Options{
				qs:     "filter[fieldA]=value1&filter[fieldB]=value2",
				Fields: []string{},
				Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"multiple filters not repeated and page",
			args{qs: "filter[fieldA]=value1&filter[fieldB]=value2&page[offset]=100"},
			Options{
				qs:     "filter[fieldA]=value1&filter[fieldB]=value2&page[offset]=100",
				Fields: []string{},
				Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
				Page:   map[string]int{"offset": 100},
				Sort:   []string{},
			},
			false,
		},
		{
			"filters and fields A",
			args{qs: "filter[fieldB]=value1&fields=fieldA,fieldB"},
			Options{
				qs:     "filter[fieldB]=value1&fields=fieldA,fieldB",
				Fields: []string{"fieldA", "fieldB"},
				Filter: map[string][]string{"fieldB": {"value1"}},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"filters and fields B",
			args{qs: "fields=fieldA,fieldB&filter[fieldB]=value1"},
			Options{
				qs:     "fields=fieldA,fieldB&filter[fieldB]=value1",
				Fields: []string{"fieldA", "fieldB"},
				Filter: map[string][]string{"fieldB": {"value1"}},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"filters and fields C",
			args{qs: "fields=fieldA,fieldB&filter[fieldB]=value1&fields=fieldC"},
			Options{
				qs:     "fields=fieldA,fieldB&filter[fieldB]=value1&fields=fieldC",
				Fields: []string{"fieldA", "fieldB", "fieldC"},
				Filter: map[string][]string{"fieldB": {"value1"}},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"single sort",
			args{qs: "sort=fieldA"},
			Options{
				qs:     "sort=fieldA",
				Fields: []string{},
				Filter: map[string][]string{},
				Page:   map[string]int{},
				Sort:   []string{"fieldA"},
			},
			false,
		},
		{
			"multiple sort parameters",
			args{qs: "sort=fieldA&sort=fieldB&sort=fieldC"},
			Options{
				qs:     "sort=fieldA&sort=fieldB&sort=fieldC",
				Fields: []string{},
				Filter: map[string][]string{},
				Page:   map[string]int{},
				Sort:   []string{"fieldA", "fieldB", "fieldC"},
			},
			false,
		},
		{
			"multiple fields via one parameter",
			args{qs: "sort=fieldA,fieldB,fieldC"},
			Options{
				qs:     "sort=fieldA,fieldB,fieldC",
				Fields: []string{},
				Filter: map[string][]string{},
				Page:   map[string]int{},
				Sort:   []string{"fieldA", "fieldB", "fieldC"},
			},
			false,
		},
		{
			"filters, pagination and sorting provided",
			args{
				ps: offsetPS,
				qs: "filter[fieldA]=value1,value2&filter[fieldB]=*test&page[offset]=10&page[limit]=10&sort=-fieldA,fieldB",
			},
			Options{
				ps:     offsetPS,
				qs:     "filter[fieldA]=value1,value2&filter[fieldB]=*test&page[offset]=10&page[limit]=10&sort=-fieldA,fieldB",
				Fields: []string{},
				Filter: map[string][]string{"fieldA": {"value1", "value2"}, "fieldB": {"*test"}},
				Page:   map[string]int{"offset": 10, "limit": 10},
				Sort:   []string{"-fieldA", "fieldB"},
			},
			false,
		},
		{
			"no filters, no sorting, but pagination",
			args{
				ps: offsetPS,
				qs: "page[limit]=10&page[offset]=10",
			},
			Options{
				ps:     offsetPS,
				qs:     "page[limit]=10&page[offset]=10",
				Fields: []string{},
				Filter: map[string][]string{},
				Page:   map[string]int{"offset": 10, "limit": 10},
				Sort:   []string{},
			},
			false,
		},
		{
			"no filters, no sorting, but fields", args{qs: "fields=fieldA,-fieldB"}, Options{
				qs:     "fields=fieldA,-fieldB",
				Fields: []string{"fieldA", "-fieldB"},
				Filter: map[string][]string{},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"no filters, no sorting, but fields in multiple params", args{qs: "fields=fieldA,-fieldB&fields=fieldC"}, Options{
				qs:     "fields=fieldA,-fieldB&fields=fieldC",
				Fields: []string{"fieldA", "-fieldB", "fieldC"},
				Filter: map[string][]string{},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"no filters, no sorting, but pagination and url decode",
			args{
				ps: offsetPS,
				qs: "page%5Blimit%5D=10&page%5Boffset%5D=10",
			}, Options{
				ps:     offsetPS,
				qs:     "page[limit]=10&page[offset]=10",
				Fields: []string{},
				Filter: map[string][]string{},
				Page:   map[string]int{"offset": 10, "limit": 10},
				Sort:   []string{},
			},
			false,
		},
		{
			"filters with lt, lte, gt and gte clauses", args{qs: "filter[iVal1]=%3C4&filter[iVal2]=%3C%3D3&filter[iVal3]=%3E1&filter[iVal4]=%3E%3D2"}, Options{
				qs:     "filter[iVal1]=<4&filter[iVal2]=<=3&filter[iVal3]=>1&filter[iVal4]=>=2",
				Fields: []string{},
				Filter: map[string][]string{
					"iVal1": {"<4"},
					"iVal2": {"<=3"},
					"iVal3": {">1"},
					"iVal4": {">=2"},
				},
				Page: map[string]int{},
				Sort: []string{},
			},
			false,
		},
		{
			"empty value for a filter",
			args{
				qs: "filter[empty]=",
			},
			Options{
				qs:     "filter[empty]=",
				Fields: []string{},
				Filter: map[string][]string{"empty": {""}},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"empty value within one of several filters",
			args{
				qs: "filter[test]=value&filter[empty]=&filter[other]=test",
			},
			Options{
				qs:     "filter[test]=value&filter[empty]=&filter[other]=test",
				Fields: []string{},
				Filter: map[string][]string{
					"empty": {""},
					"other": {"test"},
					"test":  {"value"},
				},
				Page: map[string]int{},
				Sort: []string{},
			},
			false,
		},
		{
			"improperly formatted filter in querystring",
			args{
				qs: "filter[test]=value&filter[empty]&filter[other]=test",
			},
			Options{
				qs:     "filter[test]=value&filter[empty]&filter[other]=test",
				Fields: []string{},
				Filter: map[string][]string{
					"other": {"test"},
					"test":  {"value"},
				},
				Page: map[string]int{},
				Sort: []string{},
			},
			false,
		},
		{
			"extra non-filter parameter at end with multiple filters not repeated",
			args{qs: "filter[fieldA]=value1&filter[fieldB]=value2&something=blah"},
			Options{
				qs:     "filter[fieldA]=value1&filter[fieldB]=value2&something=blah",
				Fields: []string{},
				Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
		{
			"extra non-filter parameter at beginning with multiple filters not repeated",
			args{qs: "something=blah&filter[fieldA]=value1&filter[fieldB]=value2"},
			Options{
				qs:     "something=blah&filter[fieldA]=value1&filter[fieldB]=value2",
				Fields: []string{},
				Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
				Page:   map[string]int{},
				Sort:   []string{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromQuerystring(tt.args.qs)
			if tt.args.ps != nil {
				got.SetPaginationStrategy(tt.args.ps)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("FromQuerystring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromQuerystring()\ngot:\n\t%+v\n\nwant:\n\n\t%+v", got, tt.want)
			}
		})
	}
}
