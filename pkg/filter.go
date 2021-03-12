package queryoptions

// Filter
type Filter struct {
	BeginsWith       map[string]interface{}
	Contains         map[string]interface{}
	EndsWith         map[string]interface{}
	Exact            map[string]interface{}
	GreaterThan      map[string]interface{}
	GreaterThanEqual map[string]interface{}
	LessThan         map[string]interface{}
	LessThanEqual    map[string]interface{}
	NotEqual         map[string]interface{}
}

// FilterContainer contains high level filter types
type FilterContainer struct {
	Mandatory Filter `json:"mandatory"`
	Optional  Filter `json:"optional"`
}
