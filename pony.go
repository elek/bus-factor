package main

import (
	"fmt"
	"github.com/urfave/cli"
)

func (h *Histogram) PonyDevs() []HistogramEntry {
	result := make([]HistogramEntry, 0)
	limit := h.Sum() / 2
	i := 0.0
	for _, entry := range h.SortedView() {
		i += entry.Occurrence
		result = append(result, entry)
		if i > limit {
			return result
		}
	}
	return result
}

func init() {
	App.Commands = append(App.Commands, cli.Command{
		Name: "pony",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "repo",
			},
		},
		Action: func(c *cli.Context) error {
			histogram, err := ReadGitLog(repo(c.String("repo")))
			if err != nil {
				return err
			}
			pony := histogram.PonyDevs()
			fmt.Printf("%d", len(pony))
			return nil
		},
	})
}
