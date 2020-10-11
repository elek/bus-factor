package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"strconv"
)

var App = cli.NewApp()

func main() {
	App.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "repo",
		},
		cli.BoolFlag{
			Name: "timed",
		},
	}
	App.Action = func(ctx *cli.Context) error {
		return run(repo(ctx.String("repo")), ctx.Bool("timed"))
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

func run(repo string, timed bool) error {
	var histogram Histogram
	var err error
	if timed {
		histogram, err = TimeDiscountedGitLog(repo)

	} else {
		histogram, err = ReadGitLog(repo)

	}
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
	return nil
}

