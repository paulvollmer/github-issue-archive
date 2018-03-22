VERSION=0.1.0

build:
	@go build

release:
	git tag -a v${VERSION} -m "Version ${VERSION}"
	git push origin v${VERSION}
	goreleaser

.PHONY: build release
