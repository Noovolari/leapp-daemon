package domain_alibaba

type AlibabaSessionStatus int

const (
	NotActive AlibabaSessionStatus = iota
	Pending
	Active
)
