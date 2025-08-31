package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
)

type CLI struct {
	Summary  SummaryCmd  `cmd:"" help:"Show summary report"`
	Timeline TimelineCmd `cmd:"" help:"Show timeline report"`
	Show     ShowCmd     `cmd:"" help:"Show raw grouped data"`
}

type SummaryCmd struct {
	Repo       string `help:"Repository path" short:"r" default:"."`
	ReportType string `help:"Report to show (default: all)"`
	Type       string `help:"Display only the selected report type (avl, dev, lead, pony)" short:"t"`
	Days       int    `help:"Number of days to look back" short:"d" default:"0"`
}

func (t SummaryCmd) Run() error {
	return summary(repo(t.Repo), t.Days, t.ReportType)
}

type TimelineCmd struct {
	Repo    string `help:"Repository path" short:"r" default:"."`
	Verbose bool   `help:"Verbose output" short:"v"`
	Window  string `help:"Window (Y/Q/M = year/quarter/month)" short:"w" default:"Q"`
}

func (t TimelineCmd) Run() error {
	return timeline(repo(t.Repo), t.Window, t.Verbose)
}

type ShowCmd struct {
	Repo string `help:"Repository path" short:"r" default:"."`
	Days int    `help:"Number of days to look back" short:"d" default:"0"`
}

func (t ShowCmd) Run() error {
	return show(repo(t.Repo), t.Days)
}

func main() {
	var cli CLI
	ktx := kong.Parse(&cli,
		kong.Name("bus-factor"),
		kong.Description("A git report tool."),
	)
	err := ktx.Run()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(-1)
	}
}

func repo(s string) string {
	if s != "" && s != "." {
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
		fmt.Printf("%s devp=%.2f pony=%d hhi=%f sum=%f\n", category, month.DevPower(), len(month.PonyDevs()), month.HHI(), month.Sum())
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
	fmt.Println("contribution,author")

	for _, entry := range histogram.SortedView() {
		fmt.Printf("%0.2f,%s\n", entry.Occurrence, entry.Key)
	}
	return nil
}

func summary(repo string, days int, reportType string) error {
	var histogram Histogram
	var err error
	from := time.Time{}
	if days > 0 {
		from = time.Now().Add(time.Duration(days*-1) * time.Hour * 24)
	}
	histogram, err = ReadGitLogSince(repo, from)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("Number of commits: ", histogram.Sum())
	var reports []string
	if reportType == "" {
		reports = []string{"pony", "lead", "dev", "hhi", "avl"}
	} else {
		reports = []string{reportType}
	}
	for _, r := range reports {
		switch r {
		case "pony":
			pony := histogram.PonyDevs()
			println("Pony number: " + strconv.Itoa(len(pony)))
			for _, v := range pony {
				println(fmt.Sprintf("   %s %02f", v.Key, v.Occurrence))
			}
		case "hhi":
			fmt.Printf("Herfindahl-Hirschman index: %0.2f\n", histogram.HHI())
		case "dev":
			fmt.Printf("Dev power %0.2f\n", histogram.DevPower())
		case "lead":
			fmt.Printf("Lead factor: %0.2f\n", histogram.LeadFactor())
		case "avl":
			res, err := AVL(repo)
			if err != nil {
				return err
			}
			fmt.Printf("AVL bus factor: %d\n", len(res))
			for _, author := range res {
				fmt.Println("  " + author)
			}
		}
	}
	return nil

}
