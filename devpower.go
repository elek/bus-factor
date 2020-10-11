package main

import (
	"fmt"
	"github.com/urfave/cli"
)

func (h *Histogram) DevPower() float64 {
	result := float64(0)
	base := float64(0)
	for _, entry := range h.SortedView() {
		if base == 0 {
			base = float64(entry.Occurrence)
		}
		result += float64(entry.Occurrence) / base
	}
	return result
}

func NewDoubleHistogram() DoubleHistogram {
	return DoubleHistogram{
		Histograms: make(map[string]*Histogram),
	}
}

func init() {
	App.Commands = append(App.Commands, cli.Command{
		Name: "dp",
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
			dp := histogram.DevPower()
			fmt.Printf("%f", dp)
			return nil
		},
	})
}

