package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/hatajoe/gh"
	"golang.org/x/oauth2"
)

func main() {
	var (
		r           = flag.String("r", "", "full name of repository (e.g, hatajoe/gh)")
		projectName = flag.String("p", "", "project name")
	)
	flag.Parse()

	repo := strings.Split(*r, "/")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.AccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.ProjectListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	pj, _, err := client.Repositories.ListProjects(ctx, repo[0], repo[1], opt)
	if err != nil {
		panic(err)
	}
	var allColumns []*github.ProjectColumn
	for _, p := range pj {
		if p.GetName() != *projectName {
			continue
		}
		opt := &github.ListOptions{
			PerPage: 100,
		}
		for {
			cols, res, err := client.Projects.ListProjectColumns(ctx, p.GetID(), opt)
			if err != nil {
				panic(err)
			}
			allColumns = append(allColumns, cols...)
			if res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
		break
	}
	var allCards []*github.ProjectCard
	for _, col := range allColumns {
		if col.GetName() != "Done" {
			continue
		}
		opt := &github.ListOptions{
			PerPage: 100,
		}
		for {
			cards, res, err := client.Projects.ListProjectCards(ctx, col.GetID(), opt)
			if err != nil {
				panic(err)
			}
			allCards = append(allCards, cards...)
			if res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
	}
	for _, card := range allCards {
		l := strings.Split(card.GetContentURL(), "/")
		number, err := strconv.Atoi(l[len(l)-1])
		if err != nil {
			panic(err)
		}
		issue, _, err := client.Issues.Get(ctx, repo[0], repo[1], number)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s %s\n", issue.GetHTMLURL(), issue.GetTitle())
	}
}
