package check

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"tableflip.dev/buoy/pkg/git"
	"tableflip.dev/buoy/pkg/golang"
	"tableflip.dev/buoy/pkg/needs"
)

func Check(gomod, release, domain string, ruleset git.RulesetType, verbose bool) error {
	moduleNeeds, _, err := needs.Needs([]string{gomod}, domain)
	if err != nil {
		return err
	}

	for module, packages := range moduleNeeds {
		if err := check(module, packages, release, ruleset, verbose); err != nil {
			return err
		}
	}
	return nil
}

func check(module string, packages []string, release string, ruleset git.RulesetType, verbose bool) error {
	this, err := semver.ParseTolerant(release)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("%s\n", module)
	}
	nonReady := make([]string, 0)
	for _, p := range packages {
		url := fmt.Sprintf("https://%s?go-get=1", p)
		meta, err := golang.GetMetaImport(url)
		if err != nil {
			return fmt.Errorf("unable to fetch go import %s: %v", url, err)
		}

		if meta.VCS != "git" {
			return fmt.Errorf("unknown VCS: %s", meta.VCS)
		}

		repo, err := git.GetRepo(p, meta.RepoRoot)
		if err != nil {
			return err
		}

		ref, refType := repo.BestRefFor(this, ruleset)
		switch refType {
		case git.NoRef:
			nonReady = append(nonReady, ref)
			if verbose {
				fmt.Printf("✘ %s\n", ref)
			}
		default:
			if verbose {
				fmt.Printf("✔ %s\n", ref)
			}
		}
	}

	if len(nonReady) > 0 {
		return &Error{
			Module:       module,
			Dependencies: nonReady,
		}
	}

	return nil
}

var DependencyErr = &Error{}

// Error holds the result of a failed check.
type Error struct {
	Module       string
	Dependencies []string
}

var _ error = (*Error)(nil)

func (e *Error) Is(target error) bool {
	_, is := target.(*Error)
	return is
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s not ready for release because of the following dependencies [%s]",
		e.Module,
		strings.Join(e.Dependencies, ", "))
}
