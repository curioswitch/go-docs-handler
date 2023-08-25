package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"golang.org/x/mod/modfile"
)

var projects = getProjects()

func Test() error {
	for _, project := range projects {
		if err := sh.RunV("go", "test", "-v", fmt.Sprintf("%s/...", project)); err != nil {
			return err
		}
	}
	return nil
}

func Format() error {
	if err := sh.RunV("go", "run", fmt.Sprintf("mvdan.cc/gofumpt@%s", verGoFumpt), "-l", "-w", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "run", fmt.Sprintf("github.com/rinchsan/gosimports/cmd/gosimports@%s", verGosImports), "-w",
		"-local", "github.com/curioswitch/go-docs-handler",
		"."); err != nil {
		return nil
	}
	return nil
}

func Lint() error {
	return sh.RunV("go", "run", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", verGolangCILint), "run", "--timeout", "5m")
}

// Check runs lint and tests.
func Check() {
	mg.SerialDeps(Lint, Test)
}

func FetchDocsClient() error {
	resp, err := http.Get("https://repo1.maven.org/maven2/com/linecorp/armeria/armeria/1.25.1/armeria-1.25.1.jar")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	if err := fs.WalkDir(os.DirFS("docsclient"), ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		switch path {
		case "NOTICE.txt", "web-licenses.txt", ".":
			return nil
		}
		return os.Remove(filepath.Join("docsclient", path))
	}); err != nil {
		return err
	}

	for _, f := range zipReader.File {
		path := f.Name
		if !strings.HasPrefix(path, "com/linecorp/armeria/server/docs/") || strings.HasSuffix(path, ".class") {
			continue
		}
		if f.Name[len(f.Name)-1] == '/' {
			// Directory
			continue
		}
		path = path[len("com/linecorp/armeria/server/docs/"):]
		// Windows compatibility
		if path == "assets/favicon.png" {
			path = filepath.Join("assets", "favicon.png")
		}
		file, err := f.Open()
		if err != nil {
			return err
		}
		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join("docsclient", path), content, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func getProjects() []string {
	goWork, err := os.ReadFile("go.work")
	if err != nil {
		panic(err)
	}
	work, err := modfile.ParseWork("go.work", goWork, nil)
	if err != nil {
		panic(err)
	}
	var res []string
	for _, use := range work.Use {
		res = append(res, use.Path)
	}
	return res
}
