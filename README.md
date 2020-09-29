# buoy

Given a go.mod file and a version, what are the versions this dependency should use?

## Installation

`buoy` can be installed and upgraded by running:

```shell
go get tableflip.dev/buoy
```

## Usage

```shell
Usage:
  buoy [command]

Available Commands:
  float       Find latest versions of dependencies based on a release.
  help        Help about any command

Flags:
  -h, --help   help for buoy

Use "buoy [command] --help" for more information about a command.
```

### Float

```shell
Usage:
  buoy float go.mod [flags]

Flags:
  -d, --domain string    domain filter (default "knative.dev")
  -h, --help             help for float
  -r, --release string   release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]
```

Example: 

```shell
$ buoy float $HOME/go/src/knative.dev/eventing-github/go.mod --release v0.15
knative.dev/eventing@v0.15.4
knative.dev/pkg@release-0.15
knative.dev/serving@v0.15.3
knative.dev/test-infra@release-0.15
```

Or set the domain to and target release of that dependency:

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

The goal is to find the most stable reference for a given release. Buoy will select a `ref` for a found dependency, in this order:

1. A release tag with matching major and minor; choosing the one with the highest patch version, ex: `v0.1.2`
1. If no tags, choose the release branch, ex: `release-0.1`
1. Finally, the default branch

## TODO:

- Support `go-import` with more than one import on a single page.
- Support release branch templates. For now, hardcoded to Knative style.
