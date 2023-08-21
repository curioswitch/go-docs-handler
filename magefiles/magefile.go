package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Test runs unit tests - by default, it uses wazero; set RE2_TEST_MODE=cgo or RE2_TEST_MODE=tinygo to use either, or
// RE2_TEST_EXHAUSTIVE=1 to enable exhaustive tests that may take a long time.
func Test() error {
	return sh.RunV("go", "test", "-v", "./...")
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
