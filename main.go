package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strconv"
)

var options struct {
	csvPath string
}

func init() {
	flag.StringVar(&options.csvPath, "out", "", "the file into which the size data is written in CSV format")
}

func main() {
	releases, err := listReleases()
	if err != nil {
		panic(err)
	}

	out := os.Stdout
	if options.csvPath != "" {
		out, err := os.Create(options.csvPath)
		if err != nil {
			panic(err)
		}
		defer out.Close()
	}

	cw := csv.NewWriter(out)

	for _, release := range releases {
		stats, err := collectReleaseStats(release, false)
		if err != nil {
			panic(err)
		}

		err = cw.Write([]string{
			string(release.Name),
			strconv.Itoa(stats.Normal),
			strconv.Itoa(stats.Gzip),
			strconv.Itoa(stats.Zopfli),
			strconv.Itoa(stats.Minified),
			strconv.Itoa(stats.MinGzip),
			strconv.Itoa(stats.MinZopfli),
		})
		if err != nil {
			panic(err)
		}

		cw.Flush()
	}

	for _, release := range releases {
		// jQuery versions 3.0.0+ includes a slim build
		if release.Name.Less(VersionTag("3")) {
			continue
		}

		stats, err := collectReleaseStats(release, true)
		if err != nil {
			panic(err)
		}

		err = cw.Write([]string{
			string(release.Name) + "-slim",
			strconv.Itoa(stats.Normal),
			strconv.Itoa(stats.Gzip),
			strconv.Itoa(stats.Zopfli),
			strconv.Itoa(stats.Minified),
			strconv.Itoa(stats.MinGzip),
			strconv.Itoa(stats.MinZopfli),
		})
		if err != nil {
			panic(err)
		}

		cw.Flush()
	}
}
