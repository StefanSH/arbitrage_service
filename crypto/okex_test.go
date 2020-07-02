package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOkTicker_GetData(t *testing.T) {
	okTick := &OkTicker{}
	ticks, err := okTick.GetData()
	assert.NoError(t, err)
	assert.NotEmpty(t, ticks)
}
