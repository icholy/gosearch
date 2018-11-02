package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
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

func (r Response) SortByStars() {
	sort.Slice(r.Results, func(i, j int) bool {
		return r.Results[i].Stars > r.Results[j].Stars
	})
}

func (r Response) SortByImportCount() {
	sort.Slice(r.Results, func(i, j int) bool {
		return r.Results[i].ImportCount > r.Results[j].ImportCount
	})
}

func search(q string) (*Response, error) {
	var response Response
	if err := fetchJSON(formatURL(q), &response); err != nil {
		return nil, err
	}
	for _, result := range response.Results {
		if result.Synopsis == "" {
			result.Synopsis = "<no description>"
		}
	}
	return &response, nil
}

func formatURL(query string) string {
	u, err := url.Parse("http://api.godoc.org/search")
	if err != nil {
		panic(err)
	}
	u.Query().Add("q", query)
	return u.String()
}

func fetchJSON(url string, v interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(v)
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
