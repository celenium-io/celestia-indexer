package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableName(t *testing.T) {
	constant := Constant{}
	assert.Equal(t, "constant", constant.TableName())
}
