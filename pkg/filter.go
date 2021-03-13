package queryoptions

// Filter
type Filter struct {
	BeginsWith       map[string]string
	Contains         map[string]string
	EndsWith         map[string]string
	Exact            map[string]string
	GreaterThan      map[string]string
	GreaterThanEqual map[string]string
	LessThan         map[string]string
	LessThanEqual    map[string]string
	NotEqual         map[string]string
}

// FilterContainer contains high level filter types
type FilterContainer struct {
	Mandatory Filter `json:"mandatory"`
	Optional  Filter `json:"optional"`
}
