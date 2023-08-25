package main

import (
	"fmt"
	"golang.org/x/mod/modfile"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
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
