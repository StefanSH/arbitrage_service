package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinTicker_GetData(t *testing.T) {
	binTick := &BinTicker{}
	ticks, err := binTick.GetData()
	assert.NoError(t, err)
	assert.NotEmpty(t, ticks)
}