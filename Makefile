VERSION=0.1.0

all: lint build

lint:
	@go fmt
	@golint

build:
	@go build

release:
	git tag -a v${VERSION} -m "Version ${VERSION}"
	git push origin v${VERSION}
	goreleaser

.PHONY: all lint build release
