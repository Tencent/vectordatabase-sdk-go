package entity

type AliasResult struct {
	AffectedCount int
	Collection    string
	Alias         []string
}

type SetAliasOption struct{}

type DeleteAliasOption struct{}

type DescribeAliasOption struct{}

type ListAliasOption struct{}
