package types

type Blob struct {
	Namespace    string `json:"namespace"`
	Data         string `json:"data"`
	ShareVersion int    `json:"share_version"`
	Commitment   string `json:"commitment"`
}

type Proof struct {
	Start int64    `json:"start"`
	End   int64    `json:"end"`
	Nodes []string `json:"nodes"`
}
