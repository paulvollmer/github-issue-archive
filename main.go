package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const version = "0.1.1"

var (
	ctx     = context.Background()
	client  *github.Client
	perPage = 100
	archive issueArchive

	flagOwner  string
	flagRepo   string
	flagToken  string
	flagOutput string
)

func main() {
	versionFlag := flag.Bool("v", false, "print version and exit")
	flag.StringVar(&flagOwner, "owner", "", "github owner")
	flag.StringVar(&flagRepo, "repo", "", "github repo name")
	flag.StringVar(&flagToken, "token", "", "github token")
	flag.StringVar(&flagOutput, "out", "", "file output to write")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("v%s\n", version)
		os.Exit(0)
	}
	if flagOwner == "" {
		fmt.Println("Missing github owner")
		os.Exit(1)
	}
	if flagRepo == "" {
		fmt.Println("Missing github repo name")
		os.Exit(1)
	}
	if flagToken == "" {
		fmt.Println("Missing github token")
		os.Exit(1)
	}
	if flagOutput == "" {
		flagOutput = flagOwner + "_" + flagRepo + ".json"
		fmt.Printf("Set output to %q\n", flagOutput)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: flagToken})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	getIssues(1, flagOwner, flagRepo)
	getIssuesComments(1, flagOwner, flagRepo)

	data, err := json.MarshalIndent(archive, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(flagOutput, data, 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type issueArchive struct {
	TotalIssues   int
	Issues        []*github.Issue
	TotalComments int
	Comments      []*github.IssueComment
}

func getIssues(page int, owner, repo string) {
	opt := &github.IssueListByRepoOptions{}
	opt.Page = page
	opt.PerPage = perPage
	opt.State = "all"
	issues, res, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println("NEXT", res.NextPage)
	archive.Issues = append(archive.Issues, issues...)
	archive.TotalIssues += len(issues)

	if res.NextPage > 0 {
		page++
		getIssues(page, owner, repo)
	}
}

func getIssuesComments(page int, owner, repo string) {
	optComment := &github.IssueListCommentsOptions{}
	optComment.Page = 0
	optComment.PerPage = perPage
	issueComments, _, err := client.Issues.ListComments(ctx, owner, repo, 0, optComment)
	if err != nil {
		fmt.Println(err)
	}
	archive.TotalComments = len(issueComments)
	archive.Comments = append(archive.Comments, issueComments...)
}
