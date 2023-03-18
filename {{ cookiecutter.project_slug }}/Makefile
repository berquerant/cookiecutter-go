GOMOD = go mod
GOBUILD = go build
GOTEST = go test

ROOT = $(shell git rev-parse --show-toplevel)
BIN = dist/{{ cookiecutter.project_slug }}
CMD= {{ cookiecutter.command_name }}

VERSION = $(shell git describe --tags --abbrev=0)
COMMIT = $(shell git rev-parse HEAD)
GOVERSION = $(shell go version)

LDFLAGS_PKG = main
LDFLAGS = -ldflags="-X '$(LDFLAGS_PKG).AuthorName={{ cookiecutter.author_name }}' -X '$(LDFLAGS_PKG).AuthorEmail={{ cookiecutter.author_email }}' -X '$(LDFLAGS_PKG).Version=$(VERSION)' -X '$(LDFLAGS_PKG).GoVersion=$(GOVERSION)' -X '$(LDFLAGS_PKG).Commit=$(COMMIT)' -X '$(LDFLAGS_PKG).Project={{ cookiecutter.project_slug }}' -X '$(LDFLAGS_PKG).GithubUser={{ cookiecutter.github_username }}'"

.PHONY: $(BIN)
$(BIN):
	$(GOBUILD) -v -o $@ $(LDFLAGS)

DOCKER_RUN = docker run --rm -v "$(ROOT)":/usr/src/myapp -w /usr/src/myapp
DOCKER_IMAGE = golang:{{ cookiecutter.go_version }}

.PHONY: test
test:
ifdef COOKIECUTTER_GO_DOCKER_TEST
	$(DOCKER_RUN) $(DOCKER_IMAGE) $(GOTEST) -v -cover $(LDFLAGS) ./...
else
	$(GOTEST) -v -cover $(LDFLAGS) ./...
endif

.PHONY: init
init:
	$(GOMOD) tidy
	go install golang.org/x/tools/cmd/goimports@latest
