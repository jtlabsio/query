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
		want    *Options
		wantErr bool
	}{
		{"empty querystring", args{qs: ""}, &Options{}, false},
		{"nested multiple exact filter", args{qs: "filter[mandatory][exact][fieldA]=value1&filter[mandatory][exact][fieldB]=value2"}, &Options{
			Filter: FilterContainer{
				Mandatory: Filter{
					Exact: map[string]string{"fieldA": "value1", "fieldB": "value2"},
				},
				Optional: Filter{},
			},
		}, false},
		{"nested multiple filters not repeated", args{qs: "filter[mandatory][exact][fieldA]=value1&filter[mandatory][beginsWith][fieldA]=blah"}, &Options{
			Filter: FilterContainer{
				Mandatory: Filter{
					BeginsWith: map[string]string{"fieldA": "blah"},
					Exact:      map[string]string{"fieldA": "value1"},
				},
				Optional: Filter{},
			},
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
