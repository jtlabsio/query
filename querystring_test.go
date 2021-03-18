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
			Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
			Page:   map[string]int{},
			Sort:   []string{},
		}, false},
		{"multiple filters not repeated and page", args{qs: "filter[fieldA]=value1&filter[fieldB]=value2&page[offset]=100"}, Options{
			qs:     "filter[fieldA]=value1&filter[fieldB]=value2&page[offset]=100",
			Filter: map[string][]string{"fieldA": {"value1"}, "fieldB": {"value2"}},
			Page:   map[string]int{"offset": 100},
			Sort:   []string{},
		}, false},
		{"single sort", args{qs: "sort=fieldA"}, Options{
			qs:     "sort=fieldA",
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{"fieldA"},
		}, false},
		{"multiple sort parameters", args{qs: "sort=fieldA&sort=fieldB&sort=fieldC"}, Options{
			qs:     "sort=fieldA&sort=fieldB&sort=fieldC",
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{"fieldA", "fieldB", "fieldC"},
		}, false},
		{"mulitple fields via one parameter", args{qs: "sort=fieldA,fieldB,fieldC"}, Options{
			qs:     "sort=fieldA,fieldB,fieldC",
			Filter: map[string][]string{},
			Page:   map[string]int{},
			Sort:   []string{"fieldA", "fieldB", "fieldC"},
		}, false},
		{"filters, pagination and sorting provided", args{qs: "filter[fieldA]=value1,value2&filter[fieldB]=*test&page[offset]=10&page[limit]=10&sort=-fieldA,fieldB"}, Options{
			qs: "filter[fieldA]=value1,value2&filter[fieldB]=*test&page[offset]=10&page[limit]=10&sort=-fieldA,fieldB",
			Filter: map[string][]string{"fieldA": {"value1", "value2"}, "fieldB": {"*test"}},
			Page:   map[string]int{"offset": 10, "limit": 10},
			Sort:   []string{"-fieldA", "fieldB"},
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