PACKAGES = github.com/tmpfs/pageloop github.com/tmpfs/pageloop/model github.com/tmpfs/pageloop/vdom 

bindata:
	@go-bindata -pkg pageloop data/

test:
	@go test $(PACKAGES) 

.PHONY: bindata test
