package queryoptions

import (
	"testing"
)

func TestOptions_First(t *testing.T) {
	type fields struct {
		ps     IPaginationStrategy
		qs     string
		Filter map[string][]string
		Page   map[string]int
		Sort   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"filters, sorting, but no pagination, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}, "fieldB": {"valueC"}},
			map[string]int{},
			[]string{"fieldA,-fieldB"},
		}, "filter[fieldA]=valueA,valueB&filter[fieldB]=valueC&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing limit, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0},
			[]string{},
		}, ""},
		{"no filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0, "limit": 100},
			[]string{},
		}, "page[limit]=100&page[offset]=0"},
		{"no filters, no sorting, offset not zero, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=200",
			map[string][]string{},
			map[string]int{"offset": 200, "limit": 100},
			[]string{},
		}, "page[limit]=100&page[offset]=0"},
		{"with filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 1000, "limit": 100},
			[]string{},
		}, "filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=0"},
		{"with filters, with sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 100, "limit": 100},
			[]string{"fieldA", "-fieldB"},
		}, "filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=0&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing size, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"page": 1},
			[]string{},
		}, ""},
		{"no filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[size]=100&page[page]=0",
			map[string][]string{},
			map[string]int{"size": 100, "page": 0},
			[]string{},
		}, "page[size]=100&page[page]=0"},
		{"no filters, no sorting, offset not zero, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[size]=100&page[page]=2",
			map[string][]string{},
			map[string]int{"size": 100, "page": 2},
			[]string{},
		}, "page[size]=100&page[page]=0"},
		{"with filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"page": 1000, "size": 100},
			[]string{},
		}, "filter[fieldA]=valueA,valueB&page[size]=100&page[page]=0"},
		{"with filters, with sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"size": 100, "page": 100},
			[]string{"fieldA", "-fieldB"},
		}, "filter[fieldA]=valueA,valueB&page[size]=100&page[page]=0&sort=fieldA,-fieldB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				ps:     tt.fields.ps,
				qs:     tt.fields.qs,
				Filter: tt.fields.Filter,
				Page:   tt.fields.Page,
				Sort:   tt.fields.Sort,
			}
			if got := o.First(); got != tt.want {
				t.Errorf("Options.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_Next(t *testing.T) {
	type fields struct {
		ps     IPaginationStrategy
		qs     string
		Filter map[string][]string
		Page   map[string]int
		Sort   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"filters, sorting, but no pagination, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{},
			[]string{"fieldA,-fieldB"},
		}, "filter[fieldA]=valueA,valueB&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing limit, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0},
			[]string{},
		}, ""},
		{"no filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0, "limit": 100},
			[]string{},
		}, "page[limit]=100&page[offset]=100"},
		{"no filters, no sorting, offset not zero, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=200",
			map[string][]string{},
			map[string]int{"offset": 200, "limit": 100},
			[]string{},
		}, "page[limit]=100&page[offset]=300"},
		{"with filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 1000, "limit": 100},
			[]string{},
		}, "filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=1100"},
		{"with filters, with sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 100, "limit": 100},
			[]string{"fieldA", "-fieldB"},
		}, "filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=200&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing size, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"page": 1},
			[]string{},
		}, ""},
		{"no filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[page]=0&page[size]=100",
			map[string][]string{},
			map[string]int{"size": 100, "page": 0},
			[]string{},
		}, "page[size]=100&page[page]=1"},
		{"no filters, no sorting, page not zero, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[size]=100&page[page]=2",
			map[string][]string{},
			map[string]int{"size": 100, "page": 2},
			[]string{},
		}, "page[size]=100&page[page]=3"},
		{"with filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"page": 1000, "size": 100},
			[]string{},
		}, "filter[fieldA]=valueA,valueB&page[size]=100&page[page]=1001"},
		{"with filters, with sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"size": 100, "page": 100},
			[]string{"fieldA", "-fieldB"},
		}, "filter[fieldA]=valueA,valueB&page[size]=100&page[page]=101&sort=fieldA,-fieldB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				ps:     tt.fields.ps,
				qs:     tt.fields.qs,
				Filter: tt.fields.Filter,
				Page:   tt.fields.Page,
				Sort:   tt.fields.Sort,
			}
			if got := o.Next(); got != tt.want {
				t.Errorf("Options.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_Prev(t *testing.T) {
	type fields struct {
		ps     IPaginationStrategy
		qs     string
		Filter map[string][]string
		Page   map[string]int
		Sort   []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"filters, sorting, but no pagination, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{},
			[]string{"fieldA,-fieldB"},
		}, "filter[fieldA]=valueA,valueB&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing limit, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0},
			[]string{},
		}, ""},
		{"no filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0, "limit": 100},
			[]string{},
		}, "page[limit]=100&page[offset]=0"},
		{"no filters, no sorting, offset not zero, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=200",
			map[string][]string{},
			map[string]int{"offset": 200, "limit": 100},
			[]string{},
		}, "page[limit]=100&page[offset]=100"},
		{"with filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 1000, "limit": 100},
			[]string{},
		}, "filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=900"},
		{"with filters, with sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 100, "limit": 100},
			[]string{"fieldA", "-fieldB"},
		}, "filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=0&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing size, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"page": 1},
			[]string{},
		}, ""},
		{"no filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[page]=0&page[size]=100",
			map[string][]string{},
			map[string]int{"size": 100, "page": 0},
			[]string{},
		}, "page[size]=100&page[page]=0"},
		{"no filters, no sorting, page not zero, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[size]=100&page[page]=2",
			map[string][]string{},
			map[string]int{"size": 100, "page": 2},
			[]string{},
		}, "page[size]=100&page[page]=1"},
		{"with filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"page": 1000, "size": 100},
			[]string{},
		}, "filter[fieldA]=valueA,valueB&page[size]=100&page[page]=999"},
		{"with filters, with sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"size": 100, "page": 100},
			[]string{"fieldA", "-fieldB"},
		}, "filter[fieldA]=valueA,valueB&page[size]=100&page[page]=99&sort=fieldA,-fieldB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				ps:     tt.fields.ps,
				qs:     tt.fields.qs,
				Filter: tt.fields.Filter,
				Page:   tt.fields.Page,
				Sort:   tt.fields.Sort,
			}
			if got := o.Prev(); got != tt.want {
				t.Errorf("Options.Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_Last(t *testing.T) {
	type fields struct {
		ps     IPaginationStrategy
		qs     string
		Filter map[string][]string
		Page   map[string]int
		Sort   []string
	}
	type args struct {
		total int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"filters, sorting, but no pagination, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{},
			[]string{"fieldA,-fieldB"},
		},
			args{total: 100},
			"filter[fieldA]=valueA,valueB&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing limit, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0},
			[]string{},
		},
			args{total: 1001},
			""},
		{"no filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"offset": 0, "limit": 100},
			[]string{},
		},
			args{total: 150},
			"page[limit]=100&page[offset]=100"},
		{"no filters, no sorting, offset not zero, offset strategy", fields{
			OffsetStrategy{},
			"page[limit]=100&page[offset]=200",
			map[string][]string{},
			map[string]int{"offset": 200, "limit": 100},
			[]string{},
		},
			args{total: 907},
			"page[limit]=100&page[offset]=900"},
		{"with filters, no sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 1000, "limit": 100},
			[]string{},
		},
			args{total: 99},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=0"},
		{"with filters, with sorting, offset strategy", fields{
			OffsetStrategy{},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"offset": 100, "limit": 100},
			[]string{"fieldA", "-fieldB"},
		},
			args{total: 525},
			"filter[fieldA]=valueA,valueB&page[limit]=100&page[offset]=500&sort=fieldA,-fieldB"},
		{"no filters, no sorting, missing size, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[limit]=100&page[offset]=0",
			map[string][]string{},
			map[string]int{"page": 1},
			[]string{},
		},
			args{total: 525},
			""},
		{"no filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[page]=0&page[size]=100",
			map[string][]string{},
			map[string]int{"size": 100, "page": 0},
			[]string{},
		},
			args{total: 125},
			"page[size]=100&page[page]=1"},
		{"no filters, no sorting, page not zero, pagesize strategy", fields{
			PageSizeStrategy{},
			"page[size]=100&page[page]=2",
			map[string][]string{},
			map[string]int{"size": 100, "page": 2},
			[]string{},
		},
			args{total: 1000},
			"page[size]=100&page[page]=10"},
		{"with filters, no sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=1000",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"page": 1000, "size": 100},
			[]string{},
		},
			args{total: 507},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=5"},
		{"with filters, with sorting, pagesize strategy", fields{
			PageSizeStrategy{},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=100&sort=fieldA,-fieldB",
			map[string][]string{"fieldA": {"valueA", "valueB"}},
			map[string]int{"size": 100, "page": 100},
			[]string{"fieldA", "-fieldB"},
		},
			args{total: 30020},
			"filter[fieldA]=valueA,valueB&page[size]=100&page[page]=300&sort=fieldA,-fieldB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				ps:     tt.fields.ps,
				qs:     tt.fields.qs,
				Filter: tt.fields.Filter,
				Page:   tt.fields.Page,
				Sort:   tt.fields.Sort,
			}
			if got := o.Last(tt.args.total); got != tt.want {
				t.Errorf("Options.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}
