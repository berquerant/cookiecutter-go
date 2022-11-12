GOMOD = go mod
GOBUILD = go build
GOTEST = go test

ROOT = $(shell git rev-parse --show-toplevel)
BIN = dist/{{ cookiecutter.project_slug }}

DOCKER_RUN = docker run --rm -v "$(ROOT)":/usr/src/myapp -w /usr/src/myapp
DOCKER_IMAGE = golang:{{ cookiecutter.go_version }}

VERSION = $(shell git describe --tags --abbrev=0)
COMMIT = $(shell git rev-parse HEAD)
GOVERSION = $(shell go version)

LDFLAGS_PKG = main
LDFLAGS = -ldflags="-X '$(LDFLAGS_PKG).AuthorName={{ cookiecutter.author_name }}' -X '$(LDFLAGS_PKG).AuthorEmail={{ cookiecutter.author_email }}' -X '$(LDFLAGS_PKG).Version=$(VERSION)' -X '$(LDFLAGS_PKG).GoVersion=$(GOVERSION)' -X '$(LDFLAGS_PKG).Commit=$(COMMIT)' -X '$(LDFLAGS_PKG).Project={{ cookiecutter.project_slug }}' -X '$(LDFLAGS_PKG).GithubUser={{ cookiecutter.github_username }}'"

.PHONY: $(BIN)
$(BIN): tidy
	$(GOBUILD) -v -o $@ $(LDFLAGS)
	$(BIN) -version

.PHONY: test
test: tidy
	$(GOTEST) -v -cover $(LDFLAGS) ./...

.PHONY: tidy
tidy:
	$(GOMOD) tidy

.PHONY: docker-test
docker-test:
	$(DOCKER_RUN) $(DOCKER_IMAGE) make test

.PHONY: apple-arm
apple-arm:
	$(DOCKER_RUN) -e GOOS=darwin -e GOARCH=arm64 $(DOCKER_IMAGE) make

.PHONY: apple-amd
apple-amd:
	$(DOCKER_RUN) -e GOOS=darwin -e GOARCH=amd64 $(DOCKER_IMAGE) make

.PHONY: linux-arm
linux-arm:
	$(DOCKER_RUN) -e GOOS=linux -e GOARCH=arm64 $(DOCKER_IMAGE) make

.PHONY: linux-amd
linux-amd:
	$(DOCKER_RUN) -e GOOS=linux -e GOARCH=amd64 $(DOCKER_IMAGE) make
