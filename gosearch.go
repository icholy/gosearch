package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
)

const RESULT_TEMPLATE = `
  Results from GoDoc.org
  ----------------------

{{range .Results}}  {{.Path}} ({{.Stars}} stars, {{.ImportCount}} imports)
    {{.Synopsis}}

{{end}}
`

type Result struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Synopsis    string  `json:"synopsis"`
	Stars       int     `json:"stars"`
	Score       float64 `json:"score"`
	ImportCount int     `json:"import_count"`
}

type Response struct {
	Results []*Result `json:"results"`
}

func search(q string) (*Response, error) {
	res, err := http.Get("http://api.godoc.org/search?" + url.Values{"q": {q}}.Encode())
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(res.Body)
	response := new(Response)
	if err = decoder.Decode(response); err != nil {
		return nil, err
	}
	for _, result := range response.Results {
		if result.Synopsis == "" {
			result.Synopsis = "<no description>"
		}
	}
	return response, nil
}

func init() {
	if len(os.Args) < 2 {
		log.Fatal("usage: gosearch <query>")
	}
}

func main() {
	tmpl, err := template.New("results").Parse(RESULT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	response, err := search(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(os.Stdout, response)
}
