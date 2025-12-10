package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/celenium-io/celestia-indexer/pkg/indexer/config"
	"github.com/celenium-io/celestia-indexer/pkg/types"
)

func testFromFile(t *testing.T, filename string) {
	blockFile, err := os.Open(fmt.Sprintf("../../../test/json/block_%s.json", filename))
	if err != nil {
		t.Fatal(err)
	}
	defer blockFile.Close()

	resultFile, err := os.Open(fmt.Sprintf("../../../test/json/results_%s.json", filename))
	if err != nil {
		t.Fatal(err)
	}
	defer resultFile.Close()

	var blockData types.ResultBlock
	if err := json.NewDecoder(blockFile).Decode(&blockData); err != nil {
		t.Fatal(err)
	}

	var resultsData types.ResultBlockResults
	if err := json.NewDecoder(resultFile).Decode(&resultsData); err != nil {
		t.Fatal(err)
	}

	parser := NewModule(config.Indexer{})
	if err := parser.parse(types.BlockData{
		ResultBlock:        blockData,
		ResultBlockResults: resultsData,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestProdParser(t *testing.T) {

	tests := []string{
		"8880520",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			testFromFile(t, tt)
		})
	}

}
