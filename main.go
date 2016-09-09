package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
)

var (
	port int
)

func handler(w http.ResponseWriter, r *http.Request) {
	b := &bytes.Buffer{}
	fmt.Fprint(b, "curl ")
	fmt.Fprintf(b, "-X %q ", url.QueryEscape(r.Method))

	headers := make([]string, 0, len(r.Header))
	for header := range r.Header {
		headers = append(headers, header)
	}
	sort.Strings(headers)
	for _, header := range headers {
		values := r.Header[header]
		for _, value := range values {
			fmt.Fprintf(b, "-H \"%s: %s\" ", url.QueryEscape(header), value)
		}
	}

	r.ParseForm()
	for key, values := range r.PostForm {
		for _, value := range values {
			fmt.Fprintf(b, "-d \"%s=%s\" ", url.QueryEscape(key), url.QueryEscape(value))
		}
	}

	fmt.Fprintf(b, "'%s:%d%s'", "localhost", port, r.URL)
	fmt.Println(b.String())
}

func init() {
	flag.IntVar(&port, "p", 8080, "port to listen on")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
