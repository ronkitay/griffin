
clean:
	rm -rf dist

build:
	goreleaser build --clean --snapshot --skip=post-hooks

test: #build
	echo "Testing... TBD"

release: test 
	goreleaser release --clean 

.PHONY: build release
