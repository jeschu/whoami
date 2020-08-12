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
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//goland:noinspection GoUnhandledErrorResult
		defer r.Body.Close()
		lines := new(Lines).Init()
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
		lines.Addf("%s %s%s %s (remoteAddr=%s)", r.Method, path, query, r.Proto, r.RemoteAddr)

		for _, header := range sortedHeaderKeys(r.Header) {
			lines.Addf("%s: %v", header, strings.Join(r.Header[header], ", "))
		}
		lines.AddEmpty()
		if body, err := ioutil.ReadAll(r.Body); err == nil {
			lines.AddBytes(body)
		} else {
			lines.Addf("error reading request body: %v", err)
		}
		bodyBytes := lines.Bytes()
		fmt.Print("[")
		fmt.Print(time.Now().Format(time.RFC3339Nano))
		fmt.Print("] ")
		fmt.Println(string(bodyBytes))
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Length", strconv.Itoa(len(bodyBytes)))
		_, _ = w.Write(bodyBytes)
	})
	fmt.Println("listen at 0.0.0.0:80")
	if err := http.ListenAndServe("0.0.0.0:80", nil); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func sortedHeaderKeys(header http.Header) []string {
	var keys = make([]string, 0, len(header))
	for h := range header {
		keys = append(keys, h)
	}
	sort.Strings(keys)
	return keys
}

type Lines struct {
	lines []string
}

func (lines *Lines) Init() *Lines {
	lines.lines = make([]string, 0)
	return lines
}

func (lines *Lines) Addf(format string, a ...interface{}) *Lines {
	lines.lines = append(lines.lines, fmt.Sprintf(format, a...))
	return lines
}

func (lines *Lines) AddEmpty() *Lines {
	lines.lines = append(lines.lines, "")
	return lines
}

func (lines *Lines) AddBytes(bytes []byte) *Lines {
	lines.lines = append(lines.lines, string(bytes))
	return lines
}

func (lines *Lines) String() string { return strings.Join(lines.lines, "\n") }

func (lines *Lines) Bytes() []byte { return []byte(lines.String()) }
