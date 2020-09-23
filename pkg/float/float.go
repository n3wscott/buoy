package float

import (
	"fmt"
	"github.com/blang/semver/v4"
	"github.com/n3wscott/buoy/pkg/git"
	"github.com/n3wscott/buoy/pkg/golang"
	"golang.org/x/mod/modfile"
	"io/ioutil"
	"strings"
)

func Float(gomod, release, domain string) ([]string, error) {
	b, err := ioutil.ReadFile(gomod)
	if err != nil {
		return nil, err
	}

	file, err := modfile.Parse(gomod, b, nil)
	if err != nil {
		return nil, err
	}

	packages := make([]string, 0)
	for _, r := range file.Require {
		// Look for requirements that have the prefix of domain.
		if strings.HasPrefix(r.Mod.Path, domain) {
			packages = append(packages, r.Mod.Path)
		}
	}

	this, err := semver.ParseTolerant(release)

	refs := make([]string, 0)
	for _, p := range packages {
		meta, err := golang.GetMetaImport(p)
		if err != nil {
			panic(err)
		}

		if meta.VCS != "git" {
			return nil, fmt.Errorf("unknown VCS: %s", meta.VCS)
		}

		repo, err := git.GetRepo(p, meta.RepoRoot)
		if err != nil {
			return nil, err
		}

		refs = append(refs, repo.BestRefFor(this))
	}
	return refs, nil
}
