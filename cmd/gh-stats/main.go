package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/hatajoe/gh"
	"golang.org/x/oauth2"
)

type Author string

type File struct {
	Filename  string
	FileExt   string
	Status    string
	Additions int
	Deletions int
	Changes   int
}
type Files []*File

func (l Files) Additions() int {
	total := 0
	for _, f := range l {
		total += f.Additions
	}
	return total
}

func (l Files) Deletions() int {
	total := 0
	for _, f := range l {
		total += f.Deletions
	}
	return total
}

func (l Files) Changes() int {
	total := 0
	for _, f := range l {
		total += f.Changes
	}
	return total
}

type Commit struct {
	Total     int
	Additions int
	Deletions int
	Files     Files
}
type Commits []*Commit

func (l Commits) Total() int {
	total := 0
	for _, c := range l {
		total += c.Total
	}
	return total
}

func (l Commits) Additions() int {
	total := 0
	for _, c := range l {
		total += c.Additions
	}
	return total
}

func (l Commits) Deletions() int {
	total := 0
	for _, c := range l {
		total += c.Deletions
	}
	return total
}

func (l Commits) LinesPerCommit() float64 {
	return float64(l.Total() / len(l))
}

func (l Commits) Files() Files {
	files := Files{}
	for _, c := range l {
		files = append(files, c.Files...)
	}
	return files
}

type Stats map[Author]Commits

func (s *Stats) AddAuthor(author Author) {
	if (*s)[author] != nil {
		return
	}
	(*s)[author] = []*Commit{}
}

func (s *Stats) AddCommit(author Author, commit *Commit) {
	if (*s)[author] == nil {
		panic("author is nil")
	}
	(*s)[author] = append((*s)[author], commit)
}

func main() {
	var (
		o = flag.String("o", "", "organization name")
	)
	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.AccessToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	since := time.Now().Add(-24 * 7 * time.Hour)
	stats := &Stats{}
	var allRepos []*github.Repository
	opt := &github.ListOptions{
		PerPage: 100,
	}
	for {
		repos, res, err := client.Repositories.ListByOrg(ctx, *o, &github.RepositoryListByOrgOptions{
			ListOptions: *opt,
		})
		if err != nil {
			panic(err)
		}
		allRepos = append(allRepos, repos...)
		if res.NextPage == 0 {
			break
		}
		opt.Page = res.NextPage
	}
	fmt.Printf("since: %s\n", since)
	fmt.Printf("repositories: %d\n", len(allRepos))

	var allCommits []*github.RepositoryCommit
	for _, repo := range allRepos {
		var allCommitsInRepo []*github.RepositoryCommit
		opt := &github.ListOptions{
			PerPage: 100,
		}
		for {
			commits, res, err := client.Repositories.ListCommits(ctx, *o, repo.GetName(), &github.CommitsListOptions{
				ListOptions: *opt,
				Since:       since,
			})
			if err != nil {
				log.Println(err)
				break
			}
			var details []*github.RepositoryCommit
			for _, c := range commits {
				d, _, err := client.Repositories.GetCommit(ctx, *o, repo.GetName(), c.GetSHA())
				if err != nil {
					log.Println(err)
					continue
				}
				details = append(details, d)
			}
			allCommitsInRepo = append(allCommitsInRepo, details...)
			if res.NextPage == 0 {
				break
			}
			opt.Page = res.NextPage
		}
		if len(allCommitsInRepo) != 0 {
			fmt.Printf("%s commits %d\n", repo.GetName(), len(allCommitsInRepo))
			allCommits = append(allCommits, allCommitsInRepo...)
		}
	}
	fmt.Printf("total commits: %d\n", len(allCommits))

	for _, commit := range allCommits {
		files := Files{}
		for _, file := range commit.Files {
			filenames := strings.Split(file.GetFilename(), ".")
			fileExt := "other"
			if len(filenames) > 1 {
				fileExt = filenames[len(filenames)-1]
			}
			files = append(files, &File{
				Additions: file.GetAdditions(),
				Changes:   file.GetChanges(),
				Deletions: file.GetDeletions(),
				Filename:  file.GetFilename(),
				FileExt:   fileExt,
				Status:    file.GetStatus(),
			})
		}
		s := commit.GetStats()
		author := Author(commit.GetAuthor().GetLogin())
		stats.AddAuthor(author)
		stats.AddCommit(author, &Commit{
			Total:     s.GetTotal(),
			Additions: s.GetAdditions(),
			Deletions: s.GetDeletions(),
			Files:     files,
		})
	}

	for author, commits := range *stats {
		fileCount := 0
		files := map[string]int{}
		for _, commit := range commits {
			fileCount += len(files)
			for _, file := range commit.Files {
				files[file.FileExt]++
			}
		}

		detail := `		%s: %d files`
		details := []string{}
		for ext, cnt := range files {
			details = append(details, fmt.Sprintf(detail, ext, cnt))
		}

		format := `
` + "```" + `
author: %s
commits: %d 
	LPC: %f
	total: %d
	additions: %d
	deletions: %d
files: %d 
	changes: %d
	additions: %d
	deletions: %d
	(details)
%s
` + "```" + `

`
		fmt.Printf(
			format,
			author,
			len(commits),
			commits.LinesPerCommit(),
			commits.Total(),
			commits.Additions(),
			commits.Deletions(),
			fileCount,
			commits.Files().Changes(),
			commits.Files().Additions(),
			commits.Files().Deletions(),
			strings.Join(details, "\n"),
		)
	}
}
