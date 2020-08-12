package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		lines := make([]string, 0)
		path, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			path = r.URL.Path
		}
		query, err := url.QueryUnescape(r.URL.RawQuery)
		if err != nil {
			query = r.URL.RawQuery
		}
		if len(query) == 0 {
			query = ""
		} else {
			query = "?" + query
		}
		lines = append(lines, fmt.Sprintf("%s %s%s %s (remoteAddr=%s)", r.Method, path, query, r.Proto, r.RemoteAddr))

		headers := make([]string, 0, len(r.Header))
		for header, _ := range r.Header {
			headers = append(headers, header)
		}

		sort.Strings(headers)
		for _, header := range headers {
			lines = append(lines, fmt.Sprintf("%s: %v", header, strings.Join(r.Header[header], ", ")))
		}
		lines = append(lines, "")
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			lines = append(lines, string(body))
		} else {
			lines = append(lines, fmt.Sprintf("error reading request body: %v", err))
		}
		bodyBytes := []byte(strings.Join(lines, "\n"))
		fmt.Println(string(bodyBytes))
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Length", strconv.Itoa(len(bodyBytes)))
		w.Write(bodyBytes)
	})
	fmt.Println("listen at 0.0.0.0:80")
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
