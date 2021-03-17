package queryoptions

type IPaginationStrategy interface {
	First(map[string]int) map[string]int
	Last(map[string]int, int) map[string]int
	Next(map[string]int) map[string]int
	Prev(map[string]int) map[string]int
}

type OffsetStrategy struct{}

func (os OffsetStrategy) First(c map[string]int) map[string]int {
	ret := map[string]int{}

	// read limit
	if limit, ok := c["limit"]; ok {
		ret["limit"] = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	ret["offset"] = 0

	return ret
}

func (os OffsetStrategy) Last(c map[string]int, total int) map[string]int {
	ret := map[string]int{}

	// read limit
	if limit, ok := c["limit"]; ok {
		ret["limit"] = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	offset := total / ret["limit"] * ret["limit"]
	ret["offset"] = offset

	return ret
}

func (os OffsetStrategy) Next(c map[string]int) map[string]int {
	ret := map[string]int{}

	// read limit
	if limit, ok := c["limit"]; ok {
		ret["limit"] = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	if offset, ok := c["offset"]; ok {
		ret["offset"] = offset + ret["limit"]
	} else {
		ret["offset"] = 0
	}

	return ret
}

func (os OffsetStrategy) Prev(c map[string]int) map[string]int {
	ret := map[string]int{}

	// read limit
	if limit, ok := c["limit"]; ok {
		ret["limit"] = limit
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	if offset, ok := c["offset"]; ok {
		ret["offset"] = offset - ret["limit"]
	} else {
		ret["offset"] = 0
	}

	// don't allow the offset to go below 0
	if ret["offset"] < 0 {
		ret["offset"] = 0
	}

	return ret
}

type PageSizeStrategy struct{}

func (ps PageSizeStrategy) First(c map[string]int) map[string]int {
	ret := map[string]int{}

	// read size
	if size, ok := c["size"]; ok {
		ret["size"] = size
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	ret["page"] = 0

	return ret
}

func (os PageSizeStrategy) Last(c map[string]int, total int) map[string]int {
	ret := map[string]int{}

	// read size
	if size, ok := c["size"]; ok {
		ret["size"] = size
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	page := total / ret["size"]
	ret["page"] = page

	return ret
}

func (ps PageSizeStrategy) Next(c map[string]int) map[string]int {
	ret := map[string]int{}

	// read size
	if size, ok := c["size"]; ok {
		ret["size"] = size
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	if page, ok := c["page"]; ok {
		ret["page"] = page + 1
	} else {
		ret["page"] = 0
	}

	return ret
}

func (ps PageSizeStrategy) Prev(c map[string]int) map[string]int {
	ret := map[string]int{}

	// read size
	if size, ok := c["size"]; ok {
		ret["size"] = size
	} else {
		// if limit isn't provided, return whatever was passed in
		return c
	}

	if page, ok := c["page"]; ok {
		ret["page"] = page - 1
	} else {
		ret["page"] = 0
	}

	// don't allow the page to go below 0
	if ret["page"] < 0 {
		ret["page"] = 0
	}

	return ret
}
