package main

import (
	"flag"
	"fmt"
	"os"
)

var options struct {
	csvPath     string
	htmlPath    string
	githubtoken string
	help        bool
}

func init() {
	flag.StringVar(&options.csvPath, "csv", "", "the file into which the size data is written in CSV format (or - for stdout)")
	flag.StringVar(&options.htmlPath, "html", "", "the file into which the size data is written in HTML format (or - for stdout)")

	flag.BoolVar(&options.help, "h", false, "prints this help text")

	flag.StringVar(&options.githubtoken, "ghtoken", "", "optional GitHub personal API access token")

	flag.Parse()

	if options.help {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	releases, err := listReleases(options.githubtoken)
	if err != nil {
		panic(err)
	}

	statsList := []*SizeInfo{}

	for _, release := range releases {
		fmt.Fprintf(os.Stderr, "%s\n", release.Name)
		stats, err := collectReleaseStats(release, false)
		if err != nil {
			panic(err)
		}

		statsList = append(statsList, stats)
	}

	for _, release := range releases {
		// jQuery versions 3.0.0+ includes a slim build
		if release.Name.Less(VersionTag("3")) {
			continue
		}

		fmt.Fprintf(os.Stderr, "%s-slim\n", release.Name)
		stats, err := collectReleaseStats(release, true)
		if err != nil {
			panic(err)
		}

		statsList = append(statsList, stats)
	}

	if options.csvPath != "" {
		csvOut := os.Stdout
		if options.csvPath != "-" {
			csvOut, err = os.Create(options.csvPath)
			if err != nil {
				panic(err)
			}
			defer csvOut.Close()
		}
		err := outputCSV(csvOut, statsList)
		if err != nil {
			panic(err)
		}
	}

	if options.htmlPath != "" {
		htmlOut := os.Stdout
		if options.htmlPath != "-" {
			htmlOut, err = os.Create(options.htmlPath)
			if err != nil {
				panic(err)
			}
			defer htmlOut.Close()
		}
		err := outputHTML(htmlOut, statsList)
		if err != nil {
			panic(err)
		}
	}

}
