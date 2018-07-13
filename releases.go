package main

import (
	"context"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type VersionTag string

func (a VersionTag) Less(b VersionTag) bool {
	partsA := strings.Split(string(a), ".")
	partsB := strings.Split(string(b), ".")

	min := len(partsA)
	if len(partsB) < min {
		min = len(partsB)
	}

	for index := 0; index < min; index++ {
		partA, errA := strconv.Atoi(partsA[index])
		if errA != nil {
			return true
		}

		partB, errB := strconv.Atoi(partsB[index])
		if errB != nil {
			return false
		}

		if partA != partB {
			return partA < partB
		}
	}

	return len(partsA) < len(partsB)
}

type Release struct {
	Name VersionTag
	Hash string
}
type Releases []*Release

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Less(i, j int) bool {
	return r[i].Name.Less(r[j].Name)
}

func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Only include releases that do not have
// labels for pre-release or build metadata
var includedReleaseRegexp = regexp.MustCompile(`^(\d+\.\d+(?:\.\d+)?)$`)

func includeRelease(name string) (ok bool, version VersionTag) {
	found := includedReleaseRegexp.FindStringSubmatch(name)
	if len(found) == 0 {
		return
	}

	version = VersionTag(found[1])
	if version.Less(VersionTag("1.2")) {
		return
	}

	ok = true
	return
}

func listReleases(githubtoken string) (Releases, error) {
	var releases Releases

	ctx := context.Background()

	var httpclient *http.Client
	if githubtoken != "" {
		httpclient = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubtoken},
		))
	}
	client := github.NewClient(httpclient)

	options := &github.ListOptions{
		PerPage: 50,
	}

	for {
		tags, resp, err := client.Repositories.ListTags(
			ctx,
			jQueryOwner,
			jQueryName,
			options)
		if err != nil {
			return nil, err
		}

		for _, tag := range tags {
			if ok, version := includeRelease(*tag.Name); ok {
				releases = append(releases, &Release{
					Name: version,
					Hash: tag.Commit.GetSHA(),
				})
			}
		}

		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	// Sort by version
	sort.Sort(releases)

	// Remove duplicate hashes in sequence
	releasesDedup := releases[:1]
	lastHash := releases[0].Hash
	for _, release := range releases[1:] {
		if release.Hash != lastHash {
			releasesDedup = append(releasesDedup, release)
			lastHash = release.Hash
		}
	}

	return releasesDedup, nil
}
