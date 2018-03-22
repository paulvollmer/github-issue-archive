build:
	@go build

release: build
	zip github-issue-archive github-issue-archive.zip
