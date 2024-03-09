
clean:
	rm -rf dist

build:
	goreleaser build --clean --snapshot --skip=post-hooks

release: build 
	goreleaser release --clean 

.PHONY: build release
