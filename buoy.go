package main

import (
	"fmt"
	"github.com/blang/semver/v4"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"os"
	"strings"
)

type Build struct {
	Release string // should be a major.minor version.
	Knative []string
}

func main() {
	gomod := os.Args[1]

	b, err := ioutil.ReadFile(gomod)
	if err != nil {
		panic(err)
	}

	file, err := modfile.Parse(gomod, b, nil)
	if err != nil {
		panic(err)
	}

	_ = file

	knative := make([]string, 0)
	for _, r := range file.Require {
		if strings.Contains(r.Mod.Path, "knative.dev") {
			knative = append(knative, r.Mod.Path)
		}
	}

	build := &Build{
		Release: os.Args[2],
		Knative: knative,
	}
	this, err := semver.ParseTolerant(build.Release)
	for _, kn := range build.Knative {
		meta, err := GetMetaImport(kn)
		if err != nil {
			panic(err)
		}

		if meta.VCS != "git" {
			panic(fmt.Errorf("unknown VCS: %s", meta.VCS))
		}

		repo, err := GetRepo(kn, meta.RepoRoot)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", repo.BestRefFor(this))
	}
}
