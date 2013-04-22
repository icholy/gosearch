package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
)

const RESULT_TEMPLATE = `
  Results from GoDoc.org
  ----------------------

{{range .Results}}  {{.Path}}
    {{.Synopsis}}

{{end}}
`

type Result struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis"`
}

type Response struct {
	Results []*Result `json:"results"`
}

func search(q string) (*Response, error) {
	res, err := http.Get("http://api.godoc.org/search?q=" + q)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(res.Body)
	response := new(Response)
	if err = decoder.Decode(response); err != nil {
		return nil, err
	}
	noSynopsis := "<no description>"
	for _, result := range response.Results {
		if result.Synopsis == "" {
			result.Synopsis = noSynopsis
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
