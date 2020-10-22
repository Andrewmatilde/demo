package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

// newGithubClient new client by github tokens.
func newGithubClient(token string) *github.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return github.NewClient(httpClient)
}

// fetchAllIssuesByLabel fetch all labels from github ,it could fetch about 5000*100 issues one hour.
func fetchAllIssuesByLabel(client *github.Client, owner, name, state string, labels []string) []*github.Issue {
	pageIndex := 1
	repoOptions := github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{
			Page:    pageIndex,
			PerPage: 100,
		},
		State:  state,
		Labels: labels,
	}
	var allIssues []*github.Issue
	for {
		issues, _, err :=
			client.Issues.ListByRepo(context.Background(), owner, name, &repoOptions)
		if err != nil {
			log.Fatal(err)
		}
		allIssues = append(allIssues, issues...)
		repoOptions.Page++
		if len(issues) != 0 {
			break
		}
	}
	return allIssues
}

// marshalIssues marshal issues to byte
func marshalIssues(issues []*github.Issue) []byte {
	data, err := json.Marshal(issues)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// unpackIssues write issues to file
func unpackIssues(issues []*github.Issue, filename string) {
	issuesByte := marshalIssues(issues)
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	_, err = fp.Write(issuesByte)
	if err != nil {
		log.Fatal(err)
	}
	return
}

var cliCommand = flag.String("c", "pack", "Input Your Command")
var cliOutputFileName = flag.String("f", "data.json", "Output File Name")
var cliOwner = flag.String("o", "Andrewmatilde", "owner of repo")
var cliName = flag.String("n", "demo", "name of repo")

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	client := newGithubClient(token)
	switch *cliCommand {
	case "pack":
		issues := fetchAllIssuesByLabel(client, *cliOwner, *cliName, "all", []string{})
		unpackIssues(issues, *cliOutputFileName)
		return
	default:
		return
	}
}
