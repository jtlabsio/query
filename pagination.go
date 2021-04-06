package options

import "fmt"

type IPaginationStrategy interface {
	First(map[string]int) string
	Last(map[string]int, int) string
	Next(map[string]int) string
	Prev(map[string]int) string
}

type OffsetStrategy struct{}

func (os OffsetStrategy) First(c map[string]int) string {
	var (
		o int
		l int
	)

	// read limit
	if limit, ok := c["limit"]; ok {
		l = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return ""
	}

	o = 0

	return fmt.Sprintf("page[limit]=%d&page[offset]=%d", l, o)
}

func (os OffsetStrategy) Last(c map[string]int, total int) string {
	var (
		o int
		l int
	)

	// read limit
	if limit, ok := c["limit"]; ok {
		l = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return ""
	}

	o = total / l * l

	return fmt.Sprintf("page[limit]=%d&page[offset]=%d", l, o)
}

func (os OffsetStrategy) Next(c map[string]int) string {
	var (
		o int
		l int
	)

	// read limit
	if limit, ok := c["limit"]; ok {
		l = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return ""
	}

	if offset, ok := c["offset"]; ok {
		o = offset + l
	} else {
		o = 0
	}

	return fmt.Sprintf("page[limit]=%d&page[offset]=%d", l, o)
}

func (os OffsetStrategy) Prev(c map[string]int) string {
	var (
		o int
		l int
	)

	// read limit
	if limit, ok := c["limit"]; ok {
		l = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return ""
	}

	if offset, ok := c["offset"]; ok {
		o = offset - l
	} else {
		o = 0
	}

	// don't allow the offset to go below 0
	if o < 0 {
		o = 0
	}

	return fmt.Sprintf("page[limit]=%d&page[offset]=%d", l, o)
}

type PageSizeStrategy struct{}

func (ps PageSizeStrategy) First(c map[string]int) string {
	var (
		p int
		s int
	)

	// read size
	if size, ok := c["size"]; ok {
		s = size
	} else {
		// if size isn't provided, return whatever was passed in
		return ""
	}

	p = 0

	return fmt.Sprintf("page[size]=%d&page[page]=%d", s, p)
}

func (os PageSizeStrategy) Last(c map[string]int, total int) string {
	var (
		p int
		s int
	)

	// read size
	if size, ok := c["size"]; ok {
		s = size
	} else {
		// if size isn't provided, return whatever was passed in
		return ""
	}

	p = total / s

	return fmt.Sprintf("page[size]=%d&page[page]=%d", s, p)
}

func (ps PageSizeStrategy) Next(c map[string]int) string {
	var (
		p int
		s int
	)

	// read size
	if size, ok := c["size"]; ok {
		s = size
	} else {
		// if size isn't provided, return whatever was passed in
		return ""
	}

	if page, ok := c["page"]; ok {
		p = page + 1
	} else {
		p = 0
	}

	return fmt.Sprintf("page[size]=%d&page[page]=%d", s, p)
}

func (ps PageSizeStrategy) Prev(c map[string]int) string {
	var (
		p int
		s int
	)

	// read size
	if size, ok := c["size"]; ok {
		s = size
	} else {
		// if size isn't provided, return whatever was passed in
		return ""
	}

	if page, ok := c["page"]; ok {
		p = page - 1
	} else {
		p = 0
	}

	// don't allow the page to go below 0
	if p < 0 {
		p = 0
	}

	return fmt.Sprintf("page[size]=%d&page[page]=%d", s, p)
}
