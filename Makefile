VERSION=0.1.0

all: lint build test

lint:
	@go fmt
	@golint

build:
	@go build

test:
	./github-issues-archive -v

release:
	git tag -a v${VERSION} -m "Version ${VERSION}"
	git push origin v${VERSION}
	goreleaser

.PHONY: all lint build test release
