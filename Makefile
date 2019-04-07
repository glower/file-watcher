export CGO_ENABLED?=1
REPO=bakku-app/
GO?=go

LINT_FLAGS := run -v --deadline=120s
LINTER_EXE := golangci-lint
LINTER_BIN := ./bin/$(LINTER_EXE)
TESTFLAGS := -v -cover -tags=integration -timeout 120s

$(LINTER):
#	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.15.0
#	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.15.0
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.15.0

check: gofmt lint

test:
	$(GO) test $(TESTFLAGS) ./...

lint: $(LINTER)
	$(LINTER_BIN) $(LINT_FLAGS) ./...

GFMT=find . -not \( \( -wholename "./vendor" \) -prune \) -name "*.go" | xargs gofmt -l
gofmt:
	@UNFMT=$$($(GFMT)); if [ -n "$$UNFMT" ]; then echo "gofmt needed on" $$UNFMT && exit 1; fi
