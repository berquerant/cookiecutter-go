package main

import (
{% if cookiecutter.project_category == "Code-Generator" -%}
	"bytes"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
{%- endif %}
	"os"
	"flag"
	"fmt"
	"log"
{% if cookiecutter.project_category == "Code-Generator" -%}
	"golang.org/x/tools/go/packages"
{%- endif %}
)

const usage = `{{ cookiecutter.command_name }} - {{ cookiecutter.project_description }}

Usage:

  {{ cookiecutter.command_name }} [flags]

Flags:`

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
{% if cookiecutter.project_category == "Code-Generator" -%}
	redirectToStdout := flag.Bool("stdout", false, "print result to stdout")
	output           := flag.String("output", "", "output file name; default srcdir/{{ cookiecutter.command_name }}.go")
{%- endif %}

	log.SetFlags(0)
	log.SetPrefix("{{ cookiecutter.command_name }}: ")

	flag.Usage = Usage
	flag.Parse()

{% if cookiecutter.project_category == "Code-Generator" -%}
	g := NewGenerator()
	g.parsePackage(flag.Args())

	g.printf("// Code generated by \"{{ cookiecutter.command_name }} %s\"; DO NOT EDIT.\n\n", strings.Join(os.Args[1:], " "))
	g.printf("package %s\n\n", g.pkgName)

	if err := g.Generate(); err != nil {
		log.Panicf("during generation %v", err)
	}

	w := NewResultWriter(flag.Args(), *output, *redirectToStdout)
	if err := w.writeResult(g.Bytes()); err != nil {
		log.Panicf("write file %v", err)
	}
{%- elif cookiecutter.project_category == "Command" -%}
	fmt.Println("Hello!")
{%- endif %}
}

{% if cookiecutter.project_category == "Code-Generator" -%}
func NewResultWriter(args []string, output string, redirectToStdout bool) *ResultWriter {
	return &ResultWriter{
		args:             args,
		output:           output,
		redirectToStdout: redirectToStdout,
	}
}

type ResultWriter struct {
	args            []string
	output          string
	redirectToStdout bool
}

func (r *ResultWriter) writeResult(src []byte) error {
	if r.redirectToStdout {
		return r.writeToStdout(src)
	}
	return r.writeToDestfile(src)
}

func (r *ResultWriter) writeToStdout(src []byte) error {
	f, err := os.CreateTemp("", "{{ cookiecutter.command_name }}")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(f.Name())

	if err := r.write(src, f.Name()); err != nil {
		return err
	}
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return err
	}
	if _, err := io.Copy(os.Stdout, f); err != nil {
		return err
	}
	return nil
}

func (r *ResultWriter) writeToDestfile(src []byte) error {
	return r.write(src, r.destFilename())
}

func (r *ResultWriter) write(src []byte, fileName string) error {
	if err := os.WriteFile(fileName, src, 0600); err != nil {
		return fmt.Errorf("failed to write to %s: %w", fileName, err)
	}
	if err := r.format(fileName); err != nil {
		return fmt.Errorf("failed to goimport: %w", err)
	}
	return nil
}

func (r *ResultWriter) format(targetFile string) error {
	cmd := exec.Command("go", "run", "golang.org/x/tools/cmd/goimports@v0.7.0", "-w", targetFile)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *ResultWriter) destFilename() string {
	if r.output != "" {
		return r.output
	}
	return filepath.Join(r.destDir(), "{{ cookiecutter.command_name }}_generated.go")
}

func (r *ResultWriter) destDir() string {
	if len(r.args) == 0 {
		r.args = []string{"."}
	}
	if len(r.args) == 1 && isDirectory(r.args[0]) {
		return r.args[0]
	}
	return filepath.Dir(r.args[0])
}

func isDirectory(p string) bool {
	x, err := os.Stat(p)
	if err != nil {
		log.Fatal(err)
	}
	return x.IsDir()
}

type Generator struct {
	buf     bytes.Buffer // generated code
	pkgName string
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) printf(format string, v ...any) { fmt.Fprintf(&g.buf, format, v...) }
func (g *Generator) Bytes() []byte                  { return g.buf.Bytes() }

func (g *Generator) parsePackage(patterns []string) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName,
	}, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("%d packages found", len(pkgs))
	}
	g.pkgName = pkgs[0].Name
}

// Generate codes and write them into buf.
func (g *Generator) Generate() error {
	g.printf("func Generated() {}")
	return nil
}
{%- endif %}
