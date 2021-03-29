// Package queryoptions provides a lightweight options object
// for parsing query parameters sent in a JSONAPI friendly way
// via querystring parameters
package queryoptions

import (
	"reflect"
	"testing"
)

func TestFromQuerystring(t *testing.T) {
	type args struct {
		qs string
	}
	tests := []struct {
		name    string
		args    args
		want    Options
		wantErr bool
	}{
		{"empty querystring", args{qs: ""}, Options{}, false},
		{"multiple filters not repeated", args{qs: "filter[fieldA]=value1&filter[fieldB]=value2"}, Options{
			qs:     "filter[fieldA]=value1&filter[fieldB]=value2",
			Fields: []string{},
			Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
			Page:   map[string]int{},
			Sort:   []string{},
		}, false},
		{"multiple filters not repeated and page", args{qs: "filter[fieldA]=value1&filter[fieldB]=value2&page[offset]=100"}, Options{
			qs:     "filter[fieldA]=value1&filter[fieldB]=value2&page[offset]=100",
			Fields: []string{},
			Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
			Page:   map[string]int{"offset": 100},
			Sort:   []string{},
		}, false},
		{"single sort", args{qs: "sort=fieldA"}, Options{
			qs:     "sort=fieldA",
			Fields: []string{},
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{"fieldA"},
		}, false},
		{"multiple sort parameters", args{qs: "sort=fieldA&sort=fieldB&sort=fieldC"}, Options{
			qs:     "sort=fieldA&sort=fieldB&sort=fieldC",
			Fields: []string{},
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{"fieldA", "fieldB", "fieldC"},
		}, false},
		{"mulitple fields via one parameter", args{qs: "sort=fieldA,fieldB,fieldC"}, Options{
			qs:     "sort=fieldA,fieldB,fieldC",
			Fields: []string{},
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{"fieldA", "fieldB", "fieldC"},
		}, false},
		{"filters, pagination and sorting provided", args{qs: "filter[fieldA]=value1,value2&filter[fieldB]=*test&page[offset]=10&page[limit]=10&sort=-fieldA,fieldB"}, Options{
			qs:     "filter[fieldA]=value1,value2&filter[fieldB]=*test&page[offset]=10&page[limit]=10&sort=-fieldA,fieldB",
			Fields: []string{},
			Filter: map[string][]string{"fieldA": {"value1", "value2"}, "fieldB": {"*test"}},
			Page:   map[string]int{"offset": 10, "limit": 10},
			Sort:   []string{"-fieldA", "fieldB"},
		}, false},
		{"no filters, no sorting, but pagination", args{qs: "page[limit]=10&page[offset]=10"}, Options{
			qs:     "page[limit]=10&page[offset]=10",
			Fields: []string{},
			Filter: map[string][]string{},
			Page:   map[string]int{"offset": 10, "limit": 10},
			Sort:   []string{},
		}, false},
		{"no filters, no sorting, but fields", args{qs: "fields=fieldA,-fieldB"}, Options{
			qs:     "fields=fieldA,-fieldB",
			Fields: []string{"fieldA", "-fieldB"},
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{},
		}, false},
		{"no filters, no sorting, but fields in multiple params", args{qs: "fields=fieldA,-fieldB&fields=fieldC"}, Options{
			qs:     "fields=fieldA,-fieldB&fields=fieldC",
			Fields: []string{"fieldA", "-fieldB", "fieldC"},
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{},
		}, false},
		{"no filters, no sorting, but pagination and url decode", args{qs: "page%5Blimit%5D=10&page%5Boffset%5D=10"}, Options{
			qs:     "page[limit]=10&page[offset]=10",
			Fields: []string{},
			Filter: map[string][]string{},
			Page:   map[string]int{"offset": 10, "limit": 10},
			Sort:   []string{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromQuerystring(tt.args.qs)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromQuerystring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromQuerystring() = %v, want %v", got, tt.want)
			}
		})
	}
}
