package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

func main() {
	var (
		eventName           = os.Getenv("GITHUB_EVENT_NAME")
		eventPath           = os.Getenv("GITHUB_EVENT_PATH")
		githubToken         = os.Getenv("GITHUB_TOKEN")
		ownerThenRepo       = os.Getenv("GITHUB_REPOSITORY")
		ownerThenReposplit  = strings.Split(ownerThenRepo, "/")
		ownerName, repoName = ownerThenReposplit[0], ownerThenReposplit[1]
		ctx                 = context.Background()
		ts                  = oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		)
		tc     = oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)

		repo              *github.Repository
		pullRequestNumber int
		message           string
	)
	if githubToken == "" {
		panic(errors.New("missing env GITHUB_TOKEN"))
	}
	if len(os.Args) < 2 {
		panic(errors.New("no message trying to comment on github Pull Request"))
	}
	eventBytes, err := ioutil.ReadFile(eventPath)
	if err != nil {
		panic(err)
	}
	var event github.PullRequestEvent
	err = json.Unmarshal(eventBytes, &event)
	if err != nil {
		panic(err)
	}
	repo = event.GetRepo()
	if eventName == "pull_request" && repo != nil {
		pullRequestNumber = event.GetNumber()
	} else {
		prs, _, err := client.PullRequests.List(ctx, ownerName, repoName, nil)
		if err != nil {
			panic(err)
		}
		pushHead := event.GetAfter()
		for _, pr := range prs {
			if pr.GetHead().GetSHA() == pushHead {
				pullRequestNumber = pr.GetNumber()
			}
		}
		if pullRequestNumber == 0 {
			panic(errors.New("couldn't find an open pull request for branch with head at {pushHead}"))
		}
	}
	message = strings.Join(os.Args, " ")
	comments, _, err := client.PullRequests.ListComments(ctx, ownerName, repoName, pullRequestNumber, nil)
	for _, comment := range comments {
		if comment.GetUser().GetLogin() == "github-actions[bot]" &&
			comment.GetBody() == message {
			fmt.Println("The PR already contains this message")
			return
		}
	}
	_, _, err = client.Issues.CreateComment(
		ctx,
		ownerName,
		repoName,
		pullRequestNumber,
		&github.IssueComment{
			Body: &message,
		},
	)
	if err != nil {
		panic(err)
	}
}
