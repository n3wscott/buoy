package float

import (
	"fmt"
	"tableflip.dev/buoy/pkg/needs"

	"github.com/blang/semver/v4"
	"tableflip.dev/buoy/pkg/git"
	"tableflip.dev/buoy/pkg/golang"
)

func Float(gomod, release, domain string, strict bool) ([]string, error) {
	packages, err := needs.Needs([]string{gomod}, domain)
	if err != nil {
		return nil, err
	}

	this, err := semver.ParseTolerant(release)

	refs := make([]string, 0)
	for _, p := range packages {
		url := fmt.Sprintf("https://%s?go-get=1", p)
		meta, err := golang.GetMetaImport(url)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch go import %s: %v", url, err)
		}

		if meta.VCS != "git" {
			return nil, fmt.Errorf("unknown VCS: %s", meta.VCS)
		}

		repo, err := git.GetRepo(p, meta.RepoRoot)
		if err != nil {
			return nil, err
		}

		ref, isRelease := repo.BestRefFor(this)
		if strict {
			if isRelease {
				refs = append(refs, ref)
			}
		} else {
			refs = append(refs, ref)
		}
	}
	return refs, nil
}
