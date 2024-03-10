package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableName(t *testing.T) {
	blockSignature := BlockSignature{}
	assert.Equal(t, "block_signature", blockSignature.TableName())
}
