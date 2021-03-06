package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"path"
	"strings"
)

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

	alias, err := readAlias(repo)
	if err != nil {
		return histogram, err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		histogram.Increment(normalize(c.Author.Email, alias))
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

func normalize(email string, aliases map[string]string) string {
	if alias, found := aliases[email]; found {
		return alias
	} else {
		return email
	}
}
