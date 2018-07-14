package main

import (
	"encoding/csv"
	"io"
)

func outputCSV(writer io.Writer, stats []*SizeInfo) error {
	cw := csv.NewWriter(writer)

	cw.Write([]string{
		"Version",
		"Regular",
		"Gzipped",
		"Zopflinated",
		"Minified",
		"Minified+Gzipped",
		"Minified+Zopflinated",
	})

	for _, stat := range stats {
		err := cw.Write(stat.CSVRecord())
		if err != nil {
			return err
		}

		cw.Flush()
	}

	return nil
}
