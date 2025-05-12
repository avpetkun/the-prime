package math

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestETA(t *testing.T) {
	eta := NewETA()

	eta.UpdatePercents(0)

	time.Sleep(time.Millisecond * 100)
	d := eta.UpdatePercents(10) / time.Millisecond
	require.True(t, 900 < d && d < 950, d)

	time.Sleep(time.Millisecond * 300)
	d = eta.UpdatePercents(20) / time.Millisecond
	require.True(t, 2300 < d && d < 2500, d)

	time.Sleep(time.Millisecond * 300)
	d = eta.UpdatePercents(50) / time.Millisecond
	require.True(t, 450 < d && d < 550, d)
}
