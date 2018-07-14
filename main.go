package main

import (
	"flag"
	"log"
	"os"
	"sync"
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
	logger := log.New(os.Stderr, "", 0)

	logger.Println("Downloading list of releases...")
	releases, err := listReleases(options.githubtoken)
	if err != nil {
		panic(err)
	}

	logger.Println("Downloading releases for analysis...")

	var releaseWaitGroup sync.WaitGroup
	releaseWaitGroup.Add(len(releases))

	for _, release := range releases {
		go func(release *Release) {
			defer releaseWaitGroup.Done()
			err := populateReleaseStats(release)
			if err != nil {
				panic(err)
			}
			logger.Printf("✔ %s\n", release)
		}(release)
	}

	releaseWaitGroup.Wait()

	statsList := make([]*SizeInfo, 0, len(releases))
	for _, release := range releases {
		statsList = append(statsList, release.Stats)
	}

	// Write CSV

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

	// Write HTML

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
