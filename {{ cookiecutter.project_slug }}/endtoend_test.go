package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Verify the generator is buildable and the generated code is runnable.
//
// 1. Build the generator by go build
// 2. Generate code by the generator
// 3. Run code, testdata/*.go and the generated one.
func TestEndToEnd(t *testing.T) {
	e := newExecutor(t)
	defer e.close()

	if err := run(e.cmd, "-version"); err != nil {
		t.Fatalf("%s version %v", e.cmd, err)
	}

	for i, tc := range []struct {
		title    string
		fileName string // import and use generated code
	}{
		{
			title:    "example",
			fileName: "example.go",
		},
	} {
		i := i
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			e.compileAndRun(
				t,
				i,
				tc.fileName,
			)
		})
	}
}

type executor struct {
	dir string
	cmd string
}

func newExecutor(t *testing.T) *executor {
	t.Helper()
	e := &executor{}
	e.init(t)
	return e
}

func (e *executor) init(t *testing.T) {
	t.Helper()
	dir, err := os.MkdirTemp("", "{{ cookiecutter.project_slug }}")
	if err != nil {
		t.Fatal(err)
	}
	cmd := filepath.Join(dir, "{{ cookiecutter.project_slug }}")
	// build {{ cookiecutter.project_slug }} command
	if err := run("go", "build", "-o", cmd); err != nil {
		t.Fatal(err)
	}
	e.dir = dir
	e.cmd = cmd
}

func (e *executor) compileAndRun(
	t *testing.T,
	caseNumber int,
	fileName string,
) {
	t.Helper()
	src := filepath.Join(e.dir, fileName)
	if err := copyFile(src, filepath.Join("testdata", fileName)); err != nil {
		t.Fatal(err)
	}
	outputSrc := filepath.Join(e.dir, fmt.Sprintf("{{ cookiecutter.project_slug }}%d.go", caseNumber))
	// generate code
	// please change arguments depending on the genarator
	if err := run(
		e.cmd,
		"-output", outputSrc,
	); err != nil {
		t.Fatal(err)
	}
	// run with generated code
	if err := run("go", "run", outputSrc, src); err != nil {
		t.Fatal(err)
	}
}

func (e *executor) close() {
	os.RemoveAll(e.dir)
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(to, from string) error {
	toFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFile.Close()
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	_, err = io.Copy(toFile, fromFile)
	return err
}
