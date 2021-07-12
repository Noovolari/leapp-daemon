package domain_aws

type AwsSessionStatus int

const (
	NotActive AwsSessionStatus = iota
	Pending
	Active
)
