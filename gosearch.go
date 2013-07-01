package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Result struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis"`
}

type Response struct {
	Results []*Result `json:"results"`
}

func (r *Result) name() string {
	_, name := filepath.Split(r.Path)
	return name
}

func search(q string) ([]string, error) {
	res, err := http.Get("http://api.godoc.org/search?" + url.Values{"q": {q}}.Encode())
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(res.Body)
	response := new(Response)
	if err = decoder.Decode(response); err != nil {
		return nil, err
	}
	noSynopsis := "<no description>"
	paths := make([]string, 0, 1)
	for _, result := range response.Results {
		if result.name() != q { 
			continue 
		}
		if result.Synopsis == "" {
			result.Synopsis = noSynopsis
		}
		paths = append(paths, result.Path)
	}
	return paths, nil
}

func init() {
	if len(os.Args) < 2 {
		log.Fatal("usage: gosearch <query>")
	}
}

func main() {
	paths, err := search(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(&paths); err != nil {
		log.Fatal(err)	
	}
}
