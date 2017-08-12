PACKAGES = github.com/tmpfs/pageloop github.com/tmpfs/pageloop/model github.com/tmpfs/pageloop/vdom 

bindata:
	@go-bindata -pkg pageloop -prefix data -o assets.go $(shell find ./data -type d)

test:
	@go test $(PACKAGES) 

cover:
	@go test -cover $(PACKAGES) 

coverage:
	#@go test -coverprofile=coverage.out $(PACKAGES) 
	#@go tool cover -func=coverage.out

.PHONY: bindata test
