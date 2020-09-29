package float

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/blang/semver/v4"
	"golang.org/x/mod/modfile"
	"tableflip.dev/buoy/pkg/git"
	"tableflip.dev/buoy/pkg/golang"
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
		url := fmt.Sprintf("https://%s?go-get=1", p)
		meta, err := golang.GetMetaImport(url)
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
