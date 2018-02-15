package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	owner := flag.String("owner", "", "github owner")
	repo := flag.String("repo", "", "github repo name")
	token := flag.String("token", "", "github token")
	flag.Parse()
	if *owner == "" {
		fmt.Println("Missing github owner")
		os.Exit(1)
	}
	if *repo == "" {
		fmt.Println("Missing github repo name")
		os.Exit(1)
	}
	if *token == "" {
		fmt.Println("Missing github token")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	archive := IssueArchive{}

	opt := &github.IssueListByRepoOptions{}
	opt.Page = 1
	opt.PerPage = 100
	opt.State = "all"
	issues, _, err := client.Issues.ListByRepo(ctx, *owner, *repo, opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	archive.TotalIssues = len(issues)
	archive.Issues = issues

	optComment := &github.IssueListCommentsOptions{}
	optComment.Page = 0
	optComment.PerPage = 100
	issueComments, _, err := client.Issues.ListComments(ctx, *owner, *repo, 0, optComment)
	if err != nil {
		fmt.Println(err)
	}
	archive.TotalComments = len(issueComments)
	archive.Comments = issueComments

	data, err := json.MarshalIndent(archive, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

type IssueArchive struct {
	TotalIssues   int
	Issues        []*github.Issue
	TotalComments int
	Comments      []*github.IssueComment
}
