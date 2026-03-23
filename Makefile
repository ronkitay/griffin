
clean:
	rm -rf dist

build:
	goreleaser build --clean --snapshot --skip=post-hooks

test: build
	go test ./...

release: test 
	goreleaser release --clean 

.PHONY: build release
