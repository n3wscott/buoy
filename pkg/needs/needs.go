package needs

import (
	"io/ioutil"
	"strings"

	"golang.org/x/mod/modfile"
)

func Needs(gomod, domain string) ([]string, error) {
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

	return packages, nil
}
