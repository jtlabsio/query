package queryoptions

// Options contain filtering, pagination and sorting instructions provided via
// the querystring in bracketed object notation
type Options struct {
	qs     string
	Filter map[string][]string `json:"filter"`
	Page   map[string]int      `json:"page"`
	Sort   []string            `json:"sort"`
}

func (o Options) First() string {
	return ""
}

func (o Options) Last(total int) string {
	return ""
}

func (o Options) Next() string {
	return ""
}

func (o Options) Prev() string {
	return ""
}
