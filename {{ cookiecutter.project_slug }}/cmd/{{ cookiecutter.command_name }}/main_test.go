package main_test

import (
  "os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	e := newExecutor(t)
	defer e.close()

	if err := run(e.cmd, "-h"); err != nil {
		t.Fatalf("%s help %v", e.cmd, err)
	}

{% if cookiecutter.project_category == "Code-Generator" -%}
	for _, tc := range []endToEndTestcase{
		{
			title:    "example",
			fileName: "example.go",
			source: `package main
func main() {
	Generated()
}`,
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			tc.test(t, e)
		})
	}
{%- endif %}
}

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	dir, err := os.MkdirTemp("", "{{ cookiecutter.command_name }}")
	if err != nil {
		t.Fatal(err)
	}
	cmd := filepath.Join(dir, "{{ cookiecutter.command_name }}")
	// build {{ cookiecutter.command_name }} command
	if err := run("go", "build", "-o", cmd); err != nil {
		t.Fatal(err)
	}
	e.dir = dir
	e.cmd = cmd
}

func (e *executor) close() {
	os.RemoveAll(e.dir)
}

{% if cookiecutter.project_category == "Code-Generator" -%}
type endToEndTestcase struct {
	title    string
	fileName string
	source   string
}

func (tc *endToEndTestcase) test(t *testing.T, e *executor) {
	src := filepath.Join(e.dir, tc.fileName)
	if err := os.WriteFile(src, []byte(tc.source), 0600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(src)

	outputSrc := filepath.Join(e.dir, "{{ cookiecutter.command_name }}.go")
	// generate code
	// please change arguments depending on the genarator
	if err := run(
		e.cmd,
		"-output", outputSrc,
	); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outputSrc)
	// run with generated code
	if err := run("go", "run", outputSrc, src); err != nil {
		t.Fatal(err)
	}
}
{%- endif %}
