package main

import (
	"html/template"
	"io"
	"strconv"
)

func outputHTML(writer io.Writer, stats []*SizeInfo) error {

	tmpl, err := template.New("demo.tmpl").Funcs(template.FuncMap{

		"formatBytes": func(n int) string {
			in := strconv.Itoa(n)
			out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
			if in[0] == '-' {
				in, out[0] = in[1:], '-'
			}

			for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
				out[j] = in[i]
				if i == 0 {
					return string(out)
				}
				if k++; k == 3 {
					j, k = j-1, 0
					out[j] = ','
				}
			}
		},

		"KB": func(n int) string {
			return strconv.FormatFloat(float64(n)/1024.0, 'f', 2, 64)
		},
	}).ParseFiles("./template/demo.tmpl")

	if err != nil {
		return err
	}

	// Compute deltas
	for i, stat := range stats[1:] {
		stat.Delta = &SizeInfo{
			Normal:    stat.Normal - stats[i].Normal,
			Gzip:      stat.Gzip - stats[i].Gzip,
			Zopfli:    stat.Zopfli - stats[i].Zopfli,
			Minified:  stat.Minified - stats[i].Minified,
			MinGzip:   stat.MinGzip - stats[i].MinGzip,
			MinZopfli: stat.MinZopfli - stats[i].MinZopfli,
		}
	}

	err = tmpl.Execute(writer, &struct {
		Stats []*SizeInfo
	}{
		Stats: stats,
	})
	if err != nil {
		return err
	}

	return nil
}
