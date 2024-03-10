
clean:
	rm -rf dist

build:
	goreleaser build --clean --snapshot --skip=post-hooks

test: build
	echo "Build complete - that is all the tests for now"

release: test 
	goreleaser release --clean 

.PHONY: build release
