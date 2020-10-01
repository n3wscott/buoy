# buoy

Given a go.mod file and a version, what are the versions this dependency should
use?

## Installation

`buoy` can be installed and upgraded by running:

```shell
go get tableflip.dev/buoy
```

## Usage

```
Usage:
  buoy [command]

Available Commands:
  check       Determine if this module has release branches or releases available from each dependency for a given release.
  float       Find latest versions of dependencies based on a release.
  help        Help about any command
  needs       Find dependencies based on a base import domain.

Flags:
  -h, --help   help for buoy

Use "buoy [command] --help" for more information about a command.
```

### Check

```
Determine if this module has release branches or releases available from each dependency for a given release.

Usage:
  buoy check go.mod [flags]

Flags:
  -d, --domain string    domain filter (default "knative.dev")
  -h, --help             help for check
      --release string   release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required] (default "r")
  -v, --verbose          Print verbose output.
```

- Pass/Fail of the check controlled by the exit code: 0 = check passed. 1 =
  check failed.
- Error message for failures written to `stderr`.
- Verbose output written to `stdout`.

Example,

```
$ go run . check $HOME/go/src/knative.dev/eventing-github/go.mod --release 0.18 --verbose
[exit status 0]

$ go run . check $HOME/go/src/knative.dev/eventing-github/go.mod --release 0.19 --verbose
knative.dev/eventing-github not ready for release because of the following dependencies [knative.dev/eventing@master, knative.dev/pkg@master, knative.dev/serving@master, knative.dev/test-infra@master]
[exit status 1]
```

If you need to see a more verbose output, use `--verbose`:

```
$ go run . check $HOME/go/src/knative.dev/eventing-github/go.mod --release 0.18 --verbose
knative.dev/eventing-github
✔ knative.dev/eventing@v0.18.0
✔ knative.dev/pkg@release-0.18
✔ knative.dev/serving@v0.18.0
✔ knative.dev/test-infra@release-0.18
[exit status 0]

$ go run . check $HOME/go/src/knative.dev/eventing-github/go.mod --release 0.19 --verbose
knative.dev/eventing-github
✘ knative.dev/eventing@master
✘ knative.dev/pkg@master
✘ knative.dev/serving@master
✘ knative.dev/test-infra@master
knative.dev/eventing-github not ready for release because of the following dependencies [knative.dev/eventing@master, knative.dev/pkg@master, knative.dev/serving@master, knative.dev/test-infra@master]
[exit status 1]
```

### Float

```
Usage:
  buoy float go.mod [flags]

Flags:
  -d, --domain string    domain filter (default "knative.dev")
  -h, --help             help for float
  -r, --release string   release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]
```

Example:

```
$ buoy float go.mod --release v0.15
knative.dev/eventing@v0.15.4
knative.dev/pkg@release-0.15
knative.dev/serving@v0.15.3
knative.dev/test-infra@release-0.15
```

Or set `domain` to and target release of that dependency:

```shell script
$ buoy float go.mod --release 0.18 --domain k8s.io
k8s.io/api@v0.18.10
k8s.io/apiextensions-apiserver@v0.18.10
k8s.io/apimachinery@v0.18.10
k8s.io/client-go@v0.18.10
k8s.io/code-generator@v0.18.10
k8s.io/gengo@master
k8s.io/klog@master
```

Note: the following are equivalent releases:

- `v0.1`
- `v0.1.0`
- `0.1`
- `0.1.0`

### Float Rules

The goal is to find the most stable reference for a given release. Buoy will
select a `ref` for a found dependency, in this order:

1. A release tag with matching major and minor; choosing the one with the
   highest patch version, ex: `v0.1.2`
1. If no tags, choose the release branch, ex: `release-0.1`
1. Finally, the default branch

## Needs

```
Find dependencies based on a base import domain.

Usage:
  buoy needs go.mod [flags]

Flags:
  -d, --domain string   domain filter (default "knative.dev")
      --dot             Produce a .dot file output for use with Graphviz.
  -h, --help            help for needs
```

Example,

```
$ buoy needs $HOME/go/src/knative.dev/eventing-github/go.mod
knative.dev/eventing
knative.dev/pkg
knative.dev/serving
knative.dev/test-infra
```

Or set `domain` to see a different dependency group:

```
$ buoy needs $HOME/go/src/knative.dev/eventing-github/go.mod --domain k8s.io
k8s.io/api
k8s.io/apimachinery
k8s.io/client-go
```

Or as a graph and render using [graphvis](http://www.graphviz.org/):

```
buoy needs $HOME/go/src/knative.dev/eventing/go.mod --dot | dot -Tsvg > /tmp/kn.svg
```

## TODO:

- Support `go-import` with more than one import on a single page.
- Support release branch templates. For now, hardcoded to Knative style.
