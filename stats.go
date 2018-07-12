package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/foobaz/go-zopfli/zopfli"
)

type SizeInfo struct {
	Normal    int
	Minified  int
	Gzip      int
	Zopfli    int
	MinGzip   int
	MinZopfli int
}

func collectReleaseStats(release *Release, slim bool) (*SizeInfo, error) {
	info := &SizeInfo{}

	urlNormal := jQueryCDN
	urlMinified := jQueryMinCDN
	if slim {
		urlNormal = jQuerySlimCDN
		urlMinified = jQuerySlimMinCDN
	}

	respNormal, err := http.Get(fmt.Sprintf(urlNormal, release.Name))
	if err != nil {
		return nil, err
	}
	defer respNormal.Body.Close()

	bodyNormal, err := ioutil.ReadAll(respNormal.Body)
	if err != nil {
		return nil, err
	}

	respMinified, err := http.Get(fmt.Sprintf(urlMinified, release.Name))
	if err != nil {
		return nil, err
	}
	defer respMinified.Body.Close()

	bodyMinified, err := ioutil.ReadAll(respMinified.Body)
	if err != nil {
		return nil, err
	}

	info.Normal, info.Gzip, info.Zopfli, err = collectBodyStats(bodyNormal)
	if err != nil {
		return nil, err
	}

	info.Minified, info.MinGzip, info.MinZopfli, err = collectBodyStats(bodyMinified)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func collectBodyStats(body []byte) (normal, gzipped, zopflinated int, err error) {

	// Gzip

	gzipCounter := &Counter{}
	gw, err := gzip.NewWriterLevel(gzipCounter, 6)
	if err != nil {
		return
	}

	_, err = gw.Write(body)
	if err != nil {
		return
	}

	err = gw.Close()
	if err != nil {
		return
	}

	// Zopfli

	zopfliCounter := &Counter{}
	options := zopfli.DefaultOptions()
	err = zopfli.Compress(&options, zopfli.FORMAT_GZIP, body, zopfliCounter)
	if err != nil {
		return
	}

	normal = len(body)
	gzipped = int(gzipCounter.Count())
	zopflinated = int(zopfliCounter.Count())

	return
}
