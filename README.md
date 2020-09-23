# buoy

Given a go.mod file and a version, what are the versions this dependency should use?

## Installation

`buoy` can be installed and upgraded by running:

```shell
go get github.com/n3wscott/buoy
```

## Usage

```
buoy go.mod $release
```

Example: 

```shell
$ buoy $HOME/go/src/knative.dev/eventing-github/go.mod v0.15
knative.dev/eventing@v0.15.4
knative.dev/pkg@release-0.15
knative.dev/serving@v0.15.3
knative.dev/test-infra@release-0.15
```

Note: the following are equivalent: 

- `v0.1`
- `v0.1.0`
- `0.1`
- `0.1.0`
 

## Rules

Buoy will select a `ref` for a found dependency, in this order:

1. a release, ex: `v0.1.2`
1. a release branch, ex: `release-0.1`
1. the default branch

## TODO:

- Support `go-import` with more than one import on a single page.
- Support git urls