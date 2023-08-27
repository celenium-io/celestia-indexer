package types

// Status -
type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

var availiableStatus = map[string]struct{}{
	string(StatusFailed):  {},
	string(StatusSuccess): {},
}

func IsStatus(val string) bool {
	_, ok := availiableStatus[val]
	return ok
}
