package main

import (
	"fmt"
	cli "github.com/urfave/cli"
	"os"
	"strconv"
	"time"
)

var App = cli.NewApp()

func main() {
	App.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "repo",
		},
	}
	App.Commands = []cli.Command{
		{
			Name: "summary",

			Action: func(ctx *cli.Context) error {
				return run(repo(ctx.String("repo")))
			},
		},
		{
			Name: "timeline",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "repo",
				},
				cli.BoolFlag{
					Name: "verbose",
				},
				cli.StringFlag{
					Name:  "window",
					Value: "Q",
				},
			},
			Action: func(ctx *cli.Context) error {
				return timeline(repo(ctx.String("repo")), ctx.String("window"), ctx.Bool("verbose"))
			},
		},
		{
			Name:        "show",
			Description: "show raw grouped data",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "repo",
				},
				cli.IntFlag{
					Name:  "days",
					Value: 0,
				},
			},
			Action: func(ctx *cli.Context) error {
				return show(repo(ctx.String("repo")), ctx.Int("days"))
			},
		},
	}
	err := App.Run(os.Args)

	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(-1)
	}
}

func repo(s string) string {
	if s != "" {
		return s
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}

func timeline(repo string, window string, verbose bool) error {
	histogram, err := ReadMonthlyGitLog(repo, window)
	if err != nil {
		return err
	}
	timeline, err := histogram.Timelined()
	if err != nil {
		return err
	}
	for _, category := range timeline.Keys() {
		month := timeline.Get(category)
		fmt.Printf("%s %.2f %d %f\n", category, month.DevPower(), len(month.PonyDevs()), month.Sum())
		if verbose {
			for _, v := range month.SortedView() {
				if v.Occurrence != 0 {
					fmt.Printf("   %s %.1f\n", v.Key, v.Occurrence)
				}
			}
		}
	}

	return nil
}

func show(repo string, days int) error {
	from := time.Time{}
	if days > 0 {
		from = time.Now().Add(time.Duration(days*-1) * time.Hour * 24)
	}
	histogram, err := ReadGitLogSince(repo, from)
	if err != nil {
		return err
	}
	for _, entry := range histogram.SortedView() {
		fmt.Printf("%0.2f %s\n", entry.Occurrence, entry.Key)
	}
	return nil
}

func run(repo string) error {
	var histogram Histogram
	var err error
	histogram, err = ReadGitLog(repo)
	if err != nil {
		return err
	}

	pony := histogram.PonyDevs()
	println("Pony number: " + strconv.Itoa(len(pony)))
	for _, v := range pony {
		println(fmt.Sprintf("   %s %02f", v.Key, v.Occurrence))
	}
	fmt.Printf("Dev power %0.2f\n", histogram.DevPower())

	res, err := AVL(repo)
	if err != nil {
		return err
	}
	fmt.Printf("AVL bus factor: %d\n", len(res))
	for _, author := range res {
		fmt.Println("  " + author)
	}
	fmt.Printf("Sum of points: %0.01f\n", histogram.Sum())
	for _, v := range histogram.SortedView() {
		if v.Occurrence > 0 {
			println(fmt.Sprintf("   %s %02f", v.Key, v.Occurrence))
		}

	}
	return nil
}

