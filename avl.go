package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"math"
	"os/exec"
	"strings"
)

type FileOwnership struct {
	Files map[string]map[string]bool
}

func NewFileOwnership() FileOwnership {
	return FileOwnership{
		Files: make(map[string]map[string]bool),
	}
}

func (f *FileOwnership) AddOwner(file string, owner string) {
	if _, found := f.Files[file]; !found {
		f.Files[file] = make(map[string]bool)
	}
	f.Files[file][owner] = true
}

func (f *FileOwnership) TopOwners() []string {
	h := NewHistogram()

	for _, owners := range f.Files {
		for owner, _ := range owners {
			h.Increment(owner)
		}
	}
	result := make([]string, 0)
	for _, entry := range h.SortedView() {
		result = append(result, entry.Key)
	}
	return result
}

func (f *FileOwnership) FilesWithoutOwners() int {
	res := 0
	for _, owners := range f.Files {
		k := 0
		for _, role := range owners {
			if role {
				k++
			}
		}
		if k == 0 {
			res++
		}
	}
	return res
}

func (f *FileOwnership) AllFiles() int {
	return len(f.Files)
}

func (f *FileOwnership) RemoveOwner(owner string) {
	for file, owners := range f.Files {
		if _, found := owners[owner]; found {
			f.Files[file][owner] = false
		}
	}
}

type DOARepo struct {
	Files   map[string]*HistogramWithLast
	Deleted map[string]bool
	Renames map[string]string
}

func (r *DOARepo) UpdateFile(name string, email string) {
	if _, found := r.Files[name]; !found {
		r.Files[name] = NewHistogramWithLast()
	}
	r.Files[name].Increment(email)
}

func (r *DOARepo) Normalize(name string) string {
	t := name
	for {
		des, found := r.Renames[t]
		if !found {
			return t
		} else {
			t = des
		}
	}
	return t
}

type HistogramWithLast struct {
	Histogram
	Last string
}

func (h *HistogramWithLast) Increment(key string) {
	h.Histogram.Increment(key)
	h.Last = key
}

func NewHistogramWithLast() *HistogramWithLast {
	return &HistogramWithLast{
		NewHistogram(),
		"",
	}
}
func AVL(repo string) ([]string, error) {
	codeOwners := make([]string, 0)
	result := NewDOARepo()

	command := exec.Command("git", "log", "--name-status", "--find-renames", "--pretty=format:commit %H %ae")
	command.Dir = repo

	r, _ := command.StdoutPipe()

	done := make(chan struct{})

	scanner := bufio.NewScanner(r)

	go func() {
		currentAuthor := ""
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(strings.TrimSpace(line), " ")
			if parts[0] == "commit" {
				currentAuthor = parts[2]
			} else {
				parts = strings.Split(strings.TrimSpace(line), "\t")
				if parts[0] == "D" {
					result.Deleted[result.Normalize(parts[1])] = true
				} else if parts[0] == "M" || parts[0] == "A" {
					result.UpdateFile(result.Normalize(parts[1]), currentAuthor)
				} else if strings.HasPrefix(parts[0], "R") {
					if result.Renames[parts[2]] == parts[1] {
						delete(result.Renames, parts[2])
					} else {
						result.Renames[parts[1]] = result.Normalize(parts[2])
					}
				}

			}
		}
		done <- struct{}{}
	}()
	err := command.Start()
	if err != nil {
		return codeOwners, err
	}
	<-done
	err = command.Wait()
	if err != nil {
		return codeOwners, err
	}
	owners := NewFileOwnership()
	for file, hist := range result.Files {
		sum := hist.Histogram.Sum()
		max := float64(0)
		doas := make(map[string]float64)
		for _, entry := range hist.Histogram.SortedView() {
			doa := 3.293 + 0.164*float64(entry.Occurrence) - 0.321*math.Log(float64(1+sum-entry.Occurrence))
			if entry.Key == hist.Last {
				doa += 1.098
			}
			if doa > max {
				max = doa
			}
			doas[entry.Key] = doa
		}
		for k, v := range doas {
			doas[k] = v / max
			if doas[k] > 0.5 && v > 3.293 {
				owners.AddOwner(file, k)
			}
		}

	}
	factor := 0
	for _, owner := range owners.TopOwners() {
		if owners.FilesWithoutOwners() > owners.AllFiles()/2 {
			return codeOwners, nil
		}
		owners.RemoveOwner(owner)
		codeOwners = append(codeOwners, owner)
		factor++
	}
	return codeOwners, nil
}

func NewDOARepo() DOARepo {
	return DOARepo{
		Deleted: make(map[string]bool),
		Files:   make(map[string]*HistogramWithLast),
		Renames: make(map[string]string),
	}
}

func init() {
	App.Commands = append(App.Commands, cli.Command{
		Name: "avl",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "repo",
			},
		},
		Action: func(c *cli.Context) error {
			res, err := AVL(repo(c.String("repo")))
			if err != nil {
				return err
			}
			fmt.Printf("%d", len(res))
			return nil
		},
	})
}
