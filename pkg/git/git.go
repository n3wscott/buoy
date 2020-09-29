package git

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Repo struct {
	Ref           string
	DefaultBranch string
	Tags          []string
	Branches      []string
}

func GetRepo(ref, url string) (*Repo, error) {
	repo := new(Repo)
	repo.Ref = ref

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ref := range refs {
		if ref.Name().IsTag() {
			repo.Tags = append(repo.Tags, ref.Name().Short())
		} else if ref.Name().IsBranch() {
			repo.Branches = append(repo.Branches, ref.Name().Short())
		} else if ref.Name() == "HEAD" { // Default branch.
			repo.DefaultBranch = ref.Target().Short()
		}
	}

	return repo, nil
}

func (r *Repo) BestRefFor(this semver.Version) string {
	var largest *semver.Version

	// Look for a release.
	largest = nil
	for _, t := range r.Tags {
		if sv, ok := normalizeTagVersion(t); ok {
			v, _ := semver.Make(sv)
			if v.Major == this.Major && v.Minor == this.Minor {
				if largest == nil || largest.LT(v) {
					largest = &v
				}
			}
		}
	}
	if largest != nil {
		return fmt.Sprintf("%s@%s", r.Ref, tagVersion(*largest))
	}

	// Look for a release branch.
	largest = nil
	for _, b := range r.Branches {
		if bv, ok := normalizeBranchVersion(b); ok {
			v, _ := semver.Make(bv)

			if v.Major == this.Major && v.Minor == this.Minor {
				if largest == nil || largest.LT(v) {
					largest = &v
				}
			}
		}
	}
	if largest != nil {
		return fmt.Sprintf("%s@%s", r.Ref, branchVersion(*largest))
	}

	// Return default branch.
	return fmt.Sprintf("%s@%s", r.Ref, r.DefaultBranch)
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
