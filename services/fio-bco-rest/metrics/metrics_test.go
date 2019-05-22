package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFilename(t *testing.T) {
	assert.Equal(t, "/a/b/c.e", getFilename("/a/b/c.e"))
}

func TestMetricsPersist(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entity := NewMetrics(ctx, "", time.Hour)
	delay := 1e8
	delta := 1e8

	t.Log("TimeGetToken properly times run of GetToken function")
	{
		require.Equal(t, int64(0), entity.getTokenLatency.Count())
		entity.TimeGetToken(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.getTokenLatency.Count())
		assert.InDelta(t, entity.getTokenLatency.Percentile(0.95), delay, delta)
	}

	t.Log("TimeCreateToken properly times run of CreateToken function")
	{
		require.Equal(t, int64(0), entity.createTokenLatency.Count())
		entity.TimeCreateToken(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.createTokenLatency.Count())
		assert.InDelta(t, entity.createTokenLatency.Percentile(0.95), delay, delta)
	}

	t.Log("TimeDeleteToken properly times run of DeleteToken function")
	{
		require.Equal(t, int64(0), entity.deleteTokenLatency.Count())
		entity.TimeDeleteToken(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.deleteTokenLatency.Count())
		assert.InDelta(t, entity.deleteTokenLatency.Percentile(0.95), delay, delta)
	}
}
