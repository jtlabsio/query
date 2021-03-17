package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brozeph/queryoptions"
)

func handler(w http.ResponseWriter, r *http.Request) {
	opt, err := queryoptions.FromQuerystring(r.URL.RawQuery)
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

	for _, field := range opt.Sort {
		fmt.Println("sort by field: ", field)
	}

	fmt.Fprint(w, "request complete")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
