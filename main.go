package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const version = "0.1.0"

var (
	ctx     = context.Background()
	client  *github.Client
	perPage = 100
	archive IssueArchive

	flagOwner  string
	flagRepo   string
	flagToken  string
	flagOutput string
	flagFormat string
)

func usage() {
	fmt.Printf("Usage: github-issue-archive [flags]\n\n")
	fmt.Printf("Flags:\n")
	flag.PrintDefaults()
}

func main() {
	versionFlag := flag.Bool("v", false, "print version and exit")
	flag.StringVar(&flagOwner, "owner", "", "github owner")
	flag.StringVar(&flagRepo, "repo", "", "github repo name")
	flag.StringVar(&flagToken, "token", "", "github token")
	flag.StringVar(&flagOutput, "out", "", "file output to write")
	flag.StringVar(&flagFormat, "format", "csv", "the file format can be json or csv")
	flag.Usage = usage
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
		flagOutput = flagOwner + "_" + flagRepo + "." + flagFormat
		fmt.Printf("Set output to %q\n", flagOutput)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: flagToken})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	getIssues(1, flagOwner, flagRepo)
	getIssuesComments(1, flagOwner, flagRepo)

	switch flagFormat {
	case "json":
		writeJSON()
		break
	case "csv":
		writeCSV()
		break
	default:
		fmt.Printf("format %q not supported\n", flagFormat)
		os.Exit(1)
	}

}

// IssueArchive can store all issues and issue comments
type IssueArchive struct {
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
	issueComments, res, err := client.Issues.ListComments(ctx, owner, repo, 0, optComment)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	archive.TotalComments = len(issueComments)
	archive.Comments = append(archive.Comments, issueComments...)
	if res.NextPage > 0 {
		page++
		getIssuesComments(page, owner, repo)
	}
}

func writeJSON() {
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

func writeCSV() {
	records := [][]string{
		{"Number", "State", "Title", "User", "Comments", "URL", "Created At", "Closed At", "Body"},
	}

	for i := 0; i < len(archive.Issues); i++ {
		issue := archive.Issues[i]
		number := "#" + strconv.Itoa(issue.GetNumber())
		state := issue.GetState()
		title := issue.GetTitle()
		user := issue.User.GetLogin()
		comments := strconv.Itoa(issue.GetComments())
		url := issue.GetURL()
		createdAt := issue.GetCreatedAt().String()
		closedAt := issue.GetClosedAt().String()
		body := issue.GetBody()
		records = append(records, []string{number, state, title, user, comments, url, createdAt, closedAt, body})
	}

	f, err := os.Create(flagOutput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
