package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/plugins/openapi"
)

//go:embed petstore.yml
var spec []byte

type pet struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Tag  string `json:"tag,omitempty"`
}

var nextID = 3

var petstore = map[int]pet{
	1: {
		ID:   "1",
		Name: "garfield",
		Tag:  "cat",
	},
	2: {
		ID:   "2",
		Name: "snoopy",
		Tag:  "dog",
	},
}

func main() {
	docs, err := docshandler.New(openapi.NewPlugin(spec,
		openapi.WithExampleRequests("addPet",
			pet{
				Name: "odie",
				Tag:  "dog",
			}),
	))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/pets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pIDStr := strings.TrimPrefix(r.URL.Path, "/pets/")
		if pIDStr == "" {
			http.Error(w, "missing pet id", http.StatusBadRequest)
			return
		}

		pID, err := strconv.Atoi(pIDStr)
		if err != nil {
			http.Error(w, "invalid pet id", http.StatusBadRequest)
			return
		}
		if p, ok := petstore[pID]; ok {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(p); err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
			return
		} else {
			http.Error(w, "pet not found", http.StatusNotFound)
			return
		}
	}))

	mux.Handle("/pets", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var pets []pet
			for _, p := range petstore {
				pets = append(pets, p)
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(pets); err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
		case http.MethodPost:
			var p pet
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, "failed to decode request", http.StatusBadRequest)
				return
			}

			petstore[nextID] = p
			p.ID = strconv.Itoa(nextID)
			nextID++

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(p); err != nil {
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
		}
	}))

	mux.Handle("/docs/", http.StripPrefix("/docs", docs))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
