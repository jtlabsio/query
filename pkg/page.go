package queryoptions

type Page struct {
	Limit  rune `json:"limit"`
	Offset rune `json:"offset"`
}
