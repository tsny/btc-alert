package eps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCandlestick_DiffPercent(t *testing.T) {
	candle := NewCandlestick(60000)
	candle2 := NewCandlestick(65000)

	diff := candle.DiffPercent(candle2)
	assert.Equal(t, 8, int(diff))
}

func TestCandlestick_Diff(t *testing.T) {
	candle := NewCandlestick(60000)
	candle2 := NewCandlestick(60000)

	diff := candle.DiffPercent(candle2)
	assert.Equal(t, 0, int(diff))
}
