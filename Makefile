PACKAGES = github.com/tmpfs/pageloop github.com/tmpfs/pageloop/model github.com/tmpfs/pageloop/vdom

bindata:
	@go-bindata -pkg core -prefix data -o core/assets.go $(shell find ./data -type d)

bindata-dev:
	@go-bindata -debug -pkg core -prefix data -o core/assets.go $(shell find ./data -type d)

dev: bindata-dev
	@go run bin/main.go

build: bindata/
	@cd bin && go build -o pageloop

test:
	@go test $(PACKAGES)

cli:
	@mkcli -T data -J cli/def -M cli/man -Z cli/zsh cli/pageloop.md

cover:
	@go test -cover $(PACKAGES)

coverage:
	# TODO: run over all packages
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out

.PHONY: build bindata bindata-dev test cli
