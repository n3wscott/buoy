package git

import (
	"testing"

	"github.com/blang/semver/v4"
)

// TODO: no test for GetRepo, I did not want to mock all of git api...

func Test_Repo_BestRefFor(t *testing.T) {
	repo := &Repo{
		Ref:           "ref",
		DefaultBranch: "main",
		Tags:          []string{"v0.1.0", "bar", "v0.2.0", "baz", "v0.2.1", "v0.2.2-rc.1", "v0.2.2+build", "foo"},
		Branches:      []string{"release-0.1", "bar", "release-0.2", "baz", "main", "release-0.3"},
	}

	tests := map[string]struct {
		repo    *Repo
		version semver.Version
		want    string
		release RefType
		rule    RulesetType
	}{
		"Any - v0.1": {
			repo:    repo,
			version: semver.MustParse("0.1.0"),
			want:    "ref@v0.1.0",
			release: ReleaseRef,
			rule:    AnyRule,
		},
		"Any - v0.2": {
			repo:    repo,
			version: semver.MustParse("0.2.0"),
			want:    "ref@v0.2.1",
			release: ReleaseRef,
			rule:    AnyRule,
		},
		"Any - v0.3": {
			repo:    repo,
			version: semver.MustParse("0.3.0"),
			want:    "ref@release-0.3",
			release: ReleaseBranchRef,
			rule:    AnyRule,
		},
		"Any - v0.4": {
			repo:    repo,
			version: semver.MustParse("0.4.0"),
			want:    "ref@main",
			release: DefaultBranchRef,
			rule:    AnyRule,
		},

		"ReleaseOrReleaseBranch - v0.1": {
			repo:    repo,
			version: semver.MustParse("0.1.0"),
			want:    "ref@v0.1.0",
			release: ReleaseRef,
			rule:    ReleaseOrReleaseBranchRule,
		},
		"ReleaseOrReleaseBranch - v0.2": {
			repo:    repo,
			version: semver.MustParse("0.2.0"),
			want:    "ref@v0.2.1",
			release: ReleaseRef,
			rule:    ReleaseOrReleaseBranchRule,
		},
		"ReleaseOrReleaseBranch - v0.3": {
			repo:    repo,
			version: semver.MustParse("0.3.0"),
			want:    "ref@release-0.3",
			release: ReleaseBranchRef,
			rule:    ReleaseOrReleaseBranchRule,
		},
		"ReleaseOrReleaseBranch - v0.4": {
			repo:    repo,
			version: semver.MustParse("0.4.0"),
			want:    "ref",
			release: NoRef,
			rule:    ReleaseOrReleaseBranchRule,
		},

		"Release - v0.1": {
			repo:    repo,
			version: semver.MustParse("0.1.0"),
			want:    "ref@v0.1.0",
			release: ReleaseRef,
			rule:    ReleaseRule,
		},
		"Release - v0.2": {
			repo:    repo,
			version: semver.MustParse("0.2.0"),
			want:    "ref@v0.2.1",
			release: ReleaseRef,
			rule:    ReleaseRule,
		},
		"Release - v0.3": {
			repo:    repo,
			version: semver.MustParse("0.3.0"),
			want:    "ref",
			release: NoRef,
			rule:    ReleaseRule,
		},
		"Release - v0.4": {
			repo:    repo,
			version: semver.MustParse("0.4.0"),
			want:    "ref",
			release: NoRef,
			rule:    ReleaseRule,
		},

		"ReleaseBranch - v0.1": {
			repo:    repo,
			version: semver.MustParse("0.1.0"),
			want:    "ref@release-0.1",
			release: ReleaseBranchRef,
			rule:    ReleaseBranchRule,
		},
		"ReleaseBranch - v0.2": {
			repo:    repo,
			version: semver.MustParse("0.2.0"),
			want:    "ref@release-0.2",
			release: ReleaseBranchRef,
			rule:    ReleaseBranchRule,
		},
		"ReleaseBranch - v0.3": {
			repo:    repo,
			version: semver.MustParse("0.3.0"),
			want:    "ref@release-0.3",
			release: ReleaseBranchRef,
			rule:    ReleaseBranchRule,
		},
		"ReleaseBranch - v0.4": {
			repo:    repo,
			version: semver.MustParse("0.4.0"),
			want:    "ref",
			release: NoRef,
			rule:    ReleaseBranchRule,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, release := tt.repo.BestRefFor(tt.version, tt.rule)
			if got != tt.want {
				t.Errorf("repo.BestRefFor() got ref = %v, want %v", got, tt.want)
			}
			if release != tt.release {
				t.Errorf("repo.BestRefFor() got isRelease = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeTagVersion(t *testing.T) {
	tests := map[string]struct {
		version string
		want    string
		wantOK  bool
	}{
		"v0.1.0": {
			version: "v0.1.0",
			want:    "0.1.0",
			wantOK:  true,
		},
		"v1.2.3": {
			version: "v1.2.3",
			want:    "1.2.3",
			wantOK:  true,
		},
		"notarelease": {
			version: "notarelease",
			want:    "notarelease",
			wantOK:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotOK := normalizeTagVersion(tt.version)
			if gotOK != tt.wantOK {
				t.Errorf("normalizeBranchVersion() ok = %t, wantOK %t", gotOK, tt.wantOK)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeBranchVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tagVersion(t *testing.T) {
	tests := map[string]struct {
		version semver.Version
		want    string
	}{
		"v1.2.3": {
			version: semver.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			want: "v1.2.3",
		},
		"v0.1.0": {
			version: semver.Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			want: "v0.1.0",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tagVersion(tt.version)
			if got != tt.want {
				t.Errorf("tagVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeBranchVersion(t *testing.T) {
	tests := map[string]struct {
		version string
		want    string
		wantOK  bool
	}{
		"release-0.1": {
			version: "release-0.1",
			want:    "0.1.0",
			wantOK:  true,
		},
		"release-1.2": {
			version: "release-1.2",
			want:    "1.2.0",
			wantOK:  true,
		},
		"notarelease": {
			version: "notarelease",
			want:    "notarelease",
			wantOK:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotOK := normalizeBranchVersion(tt.version)
			if gotOK != tt.wantOK {
				t.Errorf("normalizeBranchVersion() ok = %t, wantOK %t", gotOK, tt.wantOK)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeBranchVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_branchVersion(t *testing.T) {
	tests := map[string]struct {
		version semver.Version
		want    string
	}{
		"v1.2.3": {
			version: semver.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			want: "release-1.2",
		},
		"v0.1.0": {
			version: semver.Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			want: "release-0.1",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := branchVersion(tt.version)
			if got != tt.want {
				t.Errorf("branchVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
