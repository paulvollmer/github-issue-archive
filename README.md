# github-issues-archive


## Installation
```
go get github.com/paulvollmer/github-issues-archive
```


## Usage

go to https://github.com/settings/tokens and create a github token you can use to fetch data from the api. then run
```
github-issue-archive -token 123 -owner paulvollmer -repo github-issues-archive > archive.json
```
