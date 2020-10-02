package git

import "strings"

type RulesetType int

const (
	AnyRule RulesetType = iota
	ReleaseOrReleaseBranchRule
	ReleaseRule
	ReleaseBranchRule
	InvalidRule
)

func (rt RulesetType) String() string {
	return [...]string{"Any", "ReleaseOrBranch", "Release", "Branch", "Invalid"}[rt]
}

func Ruleset(rule string) RulesetType {
	switch strings.ToLower(rule) {
	case strings.ToLower(AnyRule.String()):
		return AnyRule
	case strings.ToLower(ReleaseOrReleaseBranchRule.String()):
		return ReleaseOrReleaseBranchRule
	case strings.ToLower(ReleaseRule.String()):
		return ReleaseRule
	case strings.ToLower(ReleaseBranchRule.String()):
		return ReleaseBranchRule
	default:
		return InvalidRule
	}
}

func Rulesets() []string {
	return []string{
		AnyRule.String(),
		ReleaseOrReleaseBranchRule.String(),
		ReleaseRule.String(),
		ReleaseBranchRule.String(),
	}
}
