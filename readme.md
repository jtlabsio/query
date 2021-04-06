# JSONAPI Compliant Querystring Parser

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/go.jtlabs.io/query) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jtlabsio/query/main/LICENSE) [![Coverage](http://gocover.io/_badge/github.com/jtlabsio/query)](http://gocover.io/github.com/jtlabsio/query)


This package provides [JSONAPI](https://jsonapi.org/) compliant querystring parsing. This package can be used to extract filters, pagination and sorting details from the querystring.

## Installation

```bash
go get -u go.jtlabs.io/query
```

## Usage

The query options package is designed to be used either as middleware, or in a route handler for an HTTP request. The package can be used to parse out JSONAPI style filters, pagination details and sorting instructions as supplied via the querystring.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	options "go.jtlabs.io/query"
)

func handler(w http.ResponseWriter, r *http.Request) {
	opt, err := options.FromQuerystring(r.URL.RawQuery)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// work with the options...
	for field, filter := range opt.Filter {
		for _, value := range filter {
			fmt.Println("found supplied filter (", field, ": ", value, ")")
		}
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

	fmt.Println("pagination provided with a limit of ", limit, ", and an offset of ", offset)

	for _, field := range opt.Fields {
		fmt.Println("limit response to include (or exclude) field: ", field)
	}

	for _, field := range opt.Sort {
		fmt.Println("sort by field: ", field)
	}

	fmt.Fprint(w, "request complete")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### working example

Please see [examples/usage.go](examples/usage.go) for an example of the usage. To run the example, see as follows:

```bash
go run examples/usage.go
```

The above starts a server on `:8080` for subsequent testing. In a separate terminal window, issue a request to the example server with the following:

```bash
$ curl "localhost:8080?filter[fieldA]=value1,value2&filter[fieldB]=*test&page[offset]=10&page[limit]=10&sort=-fieldA,fieldB"
```

Notice the output in the server terminal window is as follows:

```bash
$ go run examples/usage.go
found supplied filter ( fieldA :  value1 )
found supplied filter ( fieldA :  value2 )
found supplied filter ( fieldB :  *test )
pagination provided with a limit of  10 , and an offset of  10
sort by field:  -fieldA
sort by field:  fieldB
```

### queryoptions.Options

The Options struct contains properties for the provided filters, pagination details and sorting details from the querystring.

Options.Filter is a `map[string][]string`, Options.Page is a `map[string]int` and Options.Sort is a `[]string`.

```go
options := &queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{},
  Page: map[string]int{},
  Sort: []string{}
}
```

### options.Filter

JSONAPI specifications are agnostic about how filters can be provided. However, as noted in JSONAPI recommendations (<https://jsonapi.org/recommendations/#filtering>), there is an approach that many favor for clarity.

It’s recommended that servers that wish to support filtering of a resource collection based upon associations do so by allowing query parameters that combine filter with the association name. For example, the following is a request for all comments associated with a particular post:

```http
GET /comments?filter[post]=1 HTTP/1.1
```

The above would result in the following `Options`:

```go
&queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{
    "post": {"1"}
  },
  Page: map[string]int{},
  Sort: []string{}
}
```

Multiple filter values can be combined in a comma-separated list. For example:

```http
GET /comments?filter[post]=1,2 HTTP/1.1
```

The above would result in the following `Options`:

```go
&queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{
    "post": {"1", "2"}
  },
  Page: map[string]int{},
  Sort: []string{}
}
```

Furthermore, multiple filters can be applied to a single request:

```http
GET /comments?filter[post]=1,2&filter[author]=12 HTTP/1.1
```

... which results in the following `Options`:

```go
&queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{
    "post": {"1", "2"},
    "author": {"12"}
  },
  Page: map[string]int{},
  Sort: []string{}
}
```

### options.Page

JSONAPI is also agnostic regarding pagination strategies (<https://jsonapi.org/format/#fetching-pagination>), but it is noted that numerous strategies may be used (i.e. `page[number]` and `page[size]` or `page[limit]` and `page[offset]`). The queryoptions package supports any strategy the API implements.

```http
GET /comments?page[number]=2&page[size]=100 HTTP/1.1
```

The above results in the following `Options`:

```go
&queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{},
  Page: map[string]int{
		"number": 2,
		"size": 100
	},
  Sort: []string{}
}
```

Alternatively, limit and offset may be specified:

```http
GET /comments?page[limit]=20&page[offset]=12 HTTP/1.1
```

The above results in the following `Options`:

```go
&queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{},
  Page: map[string]int{
		"limit": 20,
		"offset": 12
	},
  Sort: []string{}
}
```

Ultimately, any term provided in the brackets for `sort` will be translated to an entry in the `map[string]int` struct value.

```http
GET /comments?page[whatever]=2121 HTTP/1.1
```

... is parsed as follows:

```go
&queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{},
  Page: map[string]int{
		"whatever": 2121,
	},
  Sort: []string{}
}
```

### options.Fields

In the JSONAPI specification, sparse fieldsets are supported as an array of field names: <https://jsonapi.org/format/#fetching-sparse-fieldsets>.

```http
GET /comments?fields=fieldA&fields=-fieldB HTTP/1.1
```

and this request...

```http
GET /comments?fields=fieldA,fieldB HTTP/1.1
```

Both result in the following `Options`:

```go
&queryoptions.Options{
  Fields: []string{"fieldA","-fieldB"},
  Filter: map[string][]string{},
  Page: map[string]int{},
  Sort: []string{}
}
```

### options.Sort

In the JSONAPI specification, sorting is a simple array of fields: <https://jsonapi.org/format/#fetching-sorting>.

```http
GET /comments?sort=fieldA&sort=fieldB HTTP/1.1
```

and this request...

```http
GET /comments?sort=fieldA,fieldB HTTP/1.1
```

Both result in the following `Options`:

```go
&queryoptions.Options{
  Fields: []string{},
  Filter: map[string][]string{},
  Page: map[string]int{},
  Sort: []string{"fieldA","fieldB"}
}
```
