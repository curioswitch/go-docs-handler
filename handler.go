package docshandler

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed docsclient
var docsClient embed.FS

func New(plugins ...Plugin) (http.Handler, error) {
	mux := http.NewServeMux()

	spec, err := plugins[0].GenerateSpecification()
	if err != nil {
		return nil, err
	}

	specJSON, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("serializing spec: %w", err)
	}

	schemasJSON, err := json.Marshal(generateJSONSchema(spec))
	if err != nil {
		return nil, fmt.Errorf("serializing schemas: %w", err)
	}

	mux.Handle("/specification.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(specJSON)
	}))

	mux.Handle("/schemas.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(schemasJSON)
	}))

	mux.Handle("/versions.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Consider allowing users to provide this.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))

	mux.Handle("/injected.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Allow configuring injected scripts.
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(200)
	}))

	for _, p := range plugins {
		type hasHandler interface {
			AddToHandler(handler *http.ServeMux)
		}
		if h, ok := p.(hasHandler); ok {
			h.AddToHandler(mux)
		}
	}

	docsClientFS, _ := fs.Sub(docsClient, "docsclient")
	filesHandler := http.FileServer(http.FS(docsClientFS))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "index.html"
		}
		if r.URL.Path != "/assets/favicon.png" && !strings.HasSuffix(r.URL.Path, ".ttf") {
			w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
			r.URL.Path += ".gz"
			w.Header().Set("Content-Encoding", "gzip")
		}

		filesHandler.ServeHTTP(w, r)
	}))

	return mux, nil
}
