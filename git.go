package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"path"
	"strings"
	"time"
)

type GitAuthorAlias struct {
	Aliases map[string]string
}

func ReadGitAuthorAlias(repo string) (GitAuthorAlias, error) {
	mp, err := readAlias(repo)
	if err != nil {
		return GitAuthorAlias{}, err
	} else {
		return GitAuthorAlias{mp}, nil
	}
}

func (ga *GitAuthorAlias) Normalize(author string) string {
	if alias, found := ga.Aliases[author]; found {
		return alias
	} else {
		return author
	}
}
func ReadGitLog(repo string) (Histogram, error) {
	histogram := NewHistogram()

	r, err := git.PlainOpen(repo)
	if err != nil {
		return histogram, err
	}

	ref, err := r.Head()
	if err != nil {
		return histogram, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})

	alias, err := ReadGitAuthorAlias(repo)
	if err != nil {
		return histogram, err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		histogram.Increment(alias.Normalize(c.Author.Email))
		return nil
	})
	if err != nil {
		return histogram, err
	}
	return histogram, nil
}

func ReadGitLogSince(repo string, from time.Time) (Histogram, error) {
	histogram := NewHistogram()

	r, err := git.PlainOpen(repo)
	if err != nil {
		return histogram, err
	}

	ref, err := r.Head()
	if err != nil {
		return histogram, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})

	alias, err := ReadGitAuthorAlias(repo)
	if err != nil {
		return histogram, err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.After(from) {
			histogram.Increment(alias.Normalize(c.Author.Email))
		}
		return nil
	})
	if err != nil {
		return histogram, err
	}
	return histogram, nil
}

func ReadMonthlyGitLog(repo string, window string) (DoubleHistogram, error) {
	histogram := NewDoubleHistogram()

	r, err := git.PlainOpen(repo)
	if err != nil {
		return histogram, err
	}

	ref, err := r.Head()
	if err != nil {
		return histogram, err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})

	alias, err := ReadGitAuthorAlias(repo)
	if err != nil {
		return histogram, err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		commiterDate := c.Committer.When
		category := "unknown"
		switch window {
		case "Y":
			category = fmt.Sprintf("%d", commiterDate.Year())
		case "Q":
			category = fmt.Sprintf("%d-%02d", commiterDate.Year(), (commiterDate.Month()-1)/3+1)
		case "M":
			category = fmt.Sprintf("%d-%02d", commiterDate.Year(), commiterDate.Month())
		}

		histogram.Get(category).Increment(alias.Normalize(c.Author.Email))
		return nil
	})
	if err != nil {
		return histogram, err
	}
	return histogram, nil
}
func readAlias(repoPath string) (map[string]string, error) {
	aliasFile := path.Join(repoPath, ".git", "bus-factor-alias")
	if _, err := os.Stat(aliasFile); os.IsNotExist(err) {
		return make(map[string]string), nil
	}
	return readAliasFile(aliasFile)
}

func readAliasFile(file string) (map[string]string, error) {
	alias := make(map[string]string)

	f, err := os.Open(file)
	if err != nil {
		return alias, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			parts := strings.Split(line, ",")
			alias[parts[0]] = parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return alias, nil
}

