package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/n3wscott/buoy/pkg/git"
	"github.com/n3wscott/buoy/pkg/golang"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

func main() {
	var domain string

	var buoy = &cobra.Command{
		Use:  "buoy go.mod v0.10",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			gomod := args[0]
			b, err := ioutil.ReadFile(gomod)
			if err != nil {
				return err
			}

			file, err := modfile.Parse(gomod, b, nil)
			if err != nil {
				return err
			}

			packages := make([]string, 0)
			for _, r := range file.Require {
				if strings.Contains(r.Mod.Path, domain) {
					packages = append(packages, r.Mod.Path)
				}
			}

			this, err := semver.ParseTolerant(args[1])
			for _, p := range packages {
				meta, err := golang.GetMetaImport(p)
				if err != nil {
					panic(err)
				}

				if meta.VCS != "git" {
					panic(fmt.Errorf("unknown VCS: %s", meta.VCS))
				}

				repo, err := git.GetRepo(p, meta.RepoRoot)
				if err != nil {
					panic(err)
				}

				fmt.Printf("%s\n", repo.BestRefFor(this))
			}
			return nil
		},
	}

	buoy.Flags().StringVarP(&domain, "domain", "d", "knative.dev", "domain filter")

	if err := buoy.Execute(); err != nil {
		panic(err)
	}
}
