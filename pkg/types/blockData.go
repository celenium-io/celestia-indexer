package types

import nodeTypes "github.com/dipdup-io/celestia-indexer/pkg/node/types"

type BlockData struct {
	nodeTypes.ResultBlock
	nodeTypes.ResultBlockResults
}
