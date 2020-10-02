package git

import (
	"reflect"
	"testing"
)

func TestRuleset(t *testing.T) {
	tests := map[string]struct {
		rule string
		want RulesetType
	}{
		"Any": {
			rule: "Any",
			want: AnyRule,
		},
		"ReleaseOrBranch": {
			rule: "ReleaseOrBranch",
			want: ReleaseOrReleaseBranchRule,
		},
		"Release": {
			rule: "Release",
			want: ReleaseRule,
		},
		"Branch": {
			rule: "Branch",
			want: ReleaseBranchRule,
		},
		"Invalid": {
			rule: "Invalid",
			want: InvalidRule,
		},

		"any": {
			rule: "any",
			want: AnyRule,
		},
		"releaseorbranch": {
			rule: "ReleaseOrBranch",
			want: ReleaseOrReleaseBranchRule,
		},
		"release": {
			rule: "release",
			want: ReleaseRule,
		},
		"branch": {
			rule: "Branch",
			want: ReleaseBranchRule,
		},
		"invalid": {
			rule: "invalid",
			want: InvalidRule,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := Ruleset(tt.rule); got != tt.want {
				t.Errorf("Ruleset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRulesetType_String(t *testing.T) {
	tests := map[string]struct {
		rt   RulesetType
		want string
	}{
		"Any": {
			rt:   AnyRule,
			want: "Any",
		},
		"ReleaseOrBranch": {
			rt:   ReleaseOrReleaseBranchRule,
			want: "ReleaseOrBranch",
		},
		"Release": {
			rt:   ReleaseRule,
			want: "Release",
		},
		"Branch": {
			rt:   ReleaseBranchRule,
			want: "Branch",
		},
		"Invalid": {
			rt:   InvalidRule,
			want: "Invalid",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.rt.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRulesets(t *testing.T) {
	tests := map[string]struct {
		want []string
	}{
		"Default": {
			want: []string{"Any", "ReleaseOrBranch", "Release", "Branch"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := Rulesets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rulesets() = %v, want %v", got, tt.want)
			}
		})
	}
}
