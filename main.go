package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/google/go-github/v32/github"
)

type Build struct {
	Release string // should be a major.minor version.
	Knative []string
}

var cfg = `
release: 0.19
knative:
- knative.dev/pkg
- knative.dev/eventing
- knative.dev/serving
- knative.dev/test-infra
`

func main() {
	cfg, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	build := new(Build)
	if err := yaml.Unmarshal(cfg, build); err != nil {
		panic(err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	do := make([]string, 0)

	for _, kn := range build.Knative {
		this := thisone(client, build.Release, kn)
		do = append(do, this)
	}

	fmt.Printf("FLOATING_DEPS=(\n")
	for _, d := range do {
		fmt.Printf("  \"%s\"\n", d)
	}
	fmt.Printf(")\n")
}

func thisone(client *github.Client, release, kn string) string {
	ctx := context.Background()

	this, err := semver.ParseTolerant(release)
	if err != nil {
		panic(err)
	}

	meta, err := GetMetaImport(kn)
	if err != nil {
		panic(err)
	}

	var largest *semver.Version

	if meta.VCS != "git" {
		panic(fmt.Errorf("unknown VCS: %s", meta.VCS))
	}

	org, repo := meta.OrgRepo()

	// Check for a release.

	tags, _, err := client.Repositories.ListTags(ctx, org, repo, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
		} else {
			panic(err)
		}
	}

	largest = nil
	for _, t := range tags {
		if sv, ok := normalizeTagVersion(*t.Name); ok {
			v, _ := semver.Make(sv)
			if v.Major == this.Major && v.Minor == this.Minor {
				if largest == nil || largest.LT(v) {
					largest = &v
				}
			}
		}
	}
	if largest != nil {
		//fmt.Printf("[tag] Winner: %d.%d.%d\n", largest.Major, largest.Minor, largest.Patch)
		return fmt.Sprintf("%s@%s", kn, tagVersion(*largest))
	} else {
		//fmt.Printf("[tag] No Winner\n")
	}

	// Check for a branch.

	branch, _, err := client.Repositories.ListBranches(ctx, org, repo, nil)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
		} else {
			panic(err)
		}
	}

	largest = nil
	for _, b := range branch {
		if bv, ok := normalizeBranchVersion(*b.Name); ok {
			v, err := semver.Make(bv)
			if err != nil {
				panic(err)
			}

			if v.Major == this.Major && v.Minor == this.Minor {
				if largest == nil || largest.LT(v) {
					largest = &v
				}
			}
		}
	}
	if largest != nil {
		//fmt.Printf("[branch] Winner: %d.%d.%d\n", largest.Major, largest.Minor, largest.Patch)
		return fmt.Sprintf("%s@%s", kn, branchVersion(*largest))
	} else {
		//fmt.Printf("[branch] No Winner\n")
	}

	// Get the default branch
	r, _, err := client.Repositories.Get(ctx, org, repo)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
		}
		panic(err)
	} else {
		return fmt.Sprintf("%s@%s", kn, *r.DefaultBranch)
	}
}

func normalizeTagVersion(v string) (string, bool) {
	if strings.HasPrefix(v, "v") {
		// No need to account for unicode widths.
		return v[1:], true
	}
	return v, false
}

func tagVersion(v semver.Version) string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func normalizeBranchVersion(v string) (string, bool) {
	if strings.HasPrefix(v, "release-") {
		// No need to account for unicode widths.
		return v[len("release-"):] + ".0", true
	}
	return v, false
}

func branchVersion(v semver.Version) string {
	return fmt.Sprintf("release-%d.%d", v.Major, v.Minor)
}
