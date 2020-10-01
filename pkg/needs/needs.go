package needs

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"tableflip.dev/buoy/pkg/golang"

	"golang.org/x/mod/modfile"
)

func Needs(gomod []string, domain string) ([]string, error) {
	packages := make([]string, 0)
	for _, gm := range gomod {
		pkgs, err := needs(gm, domain)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pkgs...)
	}

	return packages, nil
}
func needs(gomod string, domain string) ([]string, error) {
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

func Dot(gomods []string, domain string) (string, error) {
	dot := new(strings.Builder)
	dot.WriteString("digraph G { \n")

	for _, gomod := range gomods {
		b, err := ioutil.ReadFile(gomod)
		if err != nil {
			return "", err
		}

		file, err := modfile.Parse(gomod, b, nil)
		if err != nil {
			return "", err
		}

		if node, err := infoString(file.Module.Mod.Path); err != nil {
			return "", err
		} else {
			dot.WriteString(node)
		}

		for _, pkg := range file.Require {
			// Look for requirements that have the prefix of domain.
			if strings.HasPrefix(pkg.Mod.Path, domain) {
				if node, err := infoString(pkg.Mod.Path); err != nil {
					return "", err
				} else {
					dot.WriteString(node)
				}

				dot.WriteString(fmt.Sprintf(" %s -> %s;\n", toKey(file.Module.Mod.Path), toKey(pkg.Mod.Path)))
			}
		}
	}
	dot.WriteString("}\n")
	return dot.String(), nil
}

func infoString(pkg string) (string, error) {
	url := fmt.Sprintf("https://%s?go-get=1", pkg)
	meta, err := golang.GetMetaImport(url)
	if err != nil {
		return "", fmt.Errorf("unable to fetch go import %s: %v", url, err)
	}

	return fmt.Sprintf(" %s [label=\"%s\", URL=\"%s\", tooltip=\"%s --> %s\"];\n",
		toKey(pkg), pkg, meta.RepoRoot, pkg, meta.RepoRoot), nil
}

func toKey(pkg string) string {
	return alphaNum.ReplaceAllString(strings.ToLower(pkg), "")
}

var alphaNum *regexp.Regexp

func init() {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	alphaNum = reg
}
