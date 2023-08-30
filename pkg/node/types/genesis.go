package types

type Genesis struct {
	GenesisTime     string      `json:"genesis_time"`
	ChainID         string      `json:"chain_id"`
	InitialHeight   uint64      `json:"initial_height,string"`
	ConsensusParams interface{} `json:"consensus_params"`
	Validators      []struct {
		Address string `json:"address"`
		PubKey  struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"pub_key"`
		Power string `json:"power"`
		Name  string `json:"name"`
	} `json:"validators"`
	AppHash  string      `json:"app_hash"`
	AppState interface{} `json:"app_state"`
}
