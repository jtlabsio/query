# queryoptions

This package provides [JSONAPI](https://jsonapi.org/) compliant querystring parsing. This package can be used to extract filters, pagination and sorting details from the querystring.

## Usage

```go
func newHandler(w http.ResponseWriter, r *http.Request) {
  opt, err := queryoptions.FromQuerystring(r.URL.RawQuery)
  if err != nil {
    // unable to parse...
  }

  // work with the options...
  for _, filter := range opt.Filter {
    // each filter heare
  }

  // handle pagination
  limit := 100
  offset := 0

  if l, ok := opt.Page["limit"]; ok {
    limit = l
  }

  if o, ok := opt.Page["offset"]; ok {
    offset = o
  }

  for _, field := range opt.Sort {
    // each sort field value here
  }
}
```

### Options

The Options struct contains properties for the provided filters, pagination details and sorting details from the querystring.

Options.Filter is a `map[string]string`, Options.Page is a `map[string]int` and Options.Sort is a `[]string`.

### Filter

JSONAPI specifications are agnostic about how filters can be provided. However, as noted in JSONAPI recommendations (<https://jsonapi.org/recommendations/#filtering>), there is an approach that many favor for clarity.

It’s recommended that servers that wish to support filtering of a resource collection based upon associations do so by allowing query parameters that combine filter with the association name. For example, the following is a request for all comments associated with a particular post:

```http
GET /comments?filter[post]=1 HTTP/1.1
```

Multiple filter values can be combined in a comma-separated list. For example:

```http
GET /comments?filter[post]=1,2 HTTP/1.1
```

Furthermore, multiple filters can be applied to a single request:

```http
GET /comments?filter[post]=1,2&filter[author]=12 HTTP/1.1
```

### Page

### Sort