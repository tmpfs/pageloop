PACKAGES = github.com/tmpfs/pageloop github.com/tmpfs/pageloop/model github.com/tmpfs/pageloop/vdom 

bindata:
	@go-bindata -pkg pageloop -prefix data -o assets.go $(shell find ./data -type d)

dev: bindata
	@go run bin/main.go

test:
	@go test $(PACKAGES) 

cover:
	@go test -cover $(PACKAGES) 

coverage:
	# TODO: run over all packages
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out

.PHONY: bindata test
