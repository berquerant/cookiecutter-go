GOMOD = go mod
GOBUILD = go build -trimpath -v
GOTEST = go test -v -cover -race

ROOT = $(shell git rev-parse --show-toplevel)
BIN = dist/{{ cookiecutter.project_slug }}
CMD = "./cmd/{{ cookiecutter.command_name }}"

.PHONY: $(BIN)
$(BIN):
	$(GOBUILD) -o $@ $(CMD)

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: init
init:
	$(GOMOD) tidy -v

.PHONY: vuln
vuln:
	go tool govulncheck ./...

.PHONY: vet
vet:
	go vet ./...

# .PHONY: generate
# generate:
# 	go generate ./...

# .PHONY: clean-generated
# clean-generated:
# 	find . -name "*_generated.go" -type f -delete
