package domain_gcp

type GcpSessionStatus int

const (
	NotActive GcpSessionStatus = iota
	Active
)
