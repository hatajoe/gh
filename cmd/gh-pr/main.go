package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/hatajoe/gh"
	"golang.org/x/oauth2"
)

func main() {
	var (
		r = flag.String("r", "", "full name of repository (e.g, hatajoe/gh)")
	)
	flag.Parse()

	repo := strings.Split(*r, "/")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.AccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	since := time.Now().Add(-24 * 7 * time.Hour)
	opt := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	var allPullRequests []*github.PullRequest

FETCHING:
	for {
		pullRequests, res, err := client.PullRequests.List(ctx, repo[0], repo[1], opt)
		if err != nil {
			panic(err)
		}
		allPullRequests = append(allPullRequests, pullRequests...)
		if res.NextPage == 0 {
			break
		}
		for _, pr := range pullRequests {
			if pr.GetCreatedAt().Before(since) {
				break FETCHING
			}
		}
		opt.Page = res.NextPage
	}

	for _, pr := range allPullRequests {
		fmt.Printf("%s %s [#%d %s](%s)\n", pr.GetCreatedAt(), pr.GetUser().GetLogin(), pr.GetNumber(), pr.GetTitle(), pr.GetHTMLURL())
	}
}
