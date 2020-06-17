package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entity := NewMetrics(ctx, "/tmp", "1", time.Hour)
	delay := 1e8
	delta := 1e8

	t.Log("TimeSync properly times synchronization latency")
	{
		require.Equal(t, int64(0), entity.syncLatency.Count())
		entity.TimeSyncLatency(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.syncLatency.Count())
		assert.InDelta(t, entity.syncLatency.Percentile(0.95), delay, delta)
	}

	t.Log("TransactionImported properly marks number of accounts imported")
	{
		require.Equal(t, int64(0), entity.importedTransactions.Count())
		entity.TransactionImported()
		assert.Equal(t, int64(1), entity.importedTransactions.Count())
	}

	t.Log("TransfersImported properly marks number of accounts exported")
	{
		require.Equal(t, int64(0), entity.importedTransfers.Count())
		entity.TransfersImported(4)
		assert.Equal(t, int64(4), entity.importedTransfers.Count())
		entity.TransfersImported(6)
		assert.Equal(t, int64(10), entity.importedTransfers.Count())
	}

}
