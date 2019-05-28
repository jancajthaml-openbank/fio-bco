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

	entity := NewMetrics(ctx, "", time.Hour)
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

	t.Log("TimeImportAccount properly times run of ImportAccount function")
	{
		require.Equal(t, int64(0), entity.importAccountLatency.Count())
		entity.TimeImportAccount(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.importAccountLatency.Count())
		assert.InDelta(t, entity.importAccountLatency.Percentile(0.95), delay, delta)
	}

	t.Log("TimeExportAccount properly times run of ExportAccount function")
	{
		require.Equal(t, int64(0), entity.exportAccountLatency.Count())
		entity.TimeExportAccount(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.exportAccountLatency.Count())
		assert.InDelta(t, entity.exportAccountLatency.Percentile(0.95), delay, delta)
	}

	t.Log("TimeImportTransaction properly times run of ImportTransaction function")
	{
		require.Equal(t, int64(0), entity.importTransactionLatency.Count())
		entity.TimeImportTransaction(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.importTransactionLatency.Count())
		assert.InDelta(t, entity.importTransactionLatency.Percentile(0.95), delay, delta)
	}

	t.Log("TimeExportTransaction properly times run of ExportTransaction function")
	{
		require.Equal(t, int64(0), entity.exportTransactionLatency.Count())
		entity.TimeExportTransaction(func() {
			time.Sleep(time.Duration(delay))
		})
		assert.Equal(t, int64(1), entity.exportTransactionLatency.Count())
		assert.InDelta(t, entity.exportTransactionLatency.Percentile(0.95), delay, delta)
	}

	t.Log("ImportedAccounts properly marks number of accounts imported")
	{
		require.Equal(t, int64(0), entity.importedAccounts.Count())
		entity.ImportedAccounts(2)
		assert.Equal(t, int64(2), entity.importedAccounts.Count())
	}

	t.Log("ExportedAccounts properly marks number of accounts exported")
	{
		require.Equal(t, int64(0), entity.exportedAccounts.Count())
		entity.ExportedAccounts(4)
		assert.Equal(t, int64(4), entity.exportedAccounts.Count())
	}

	t.Log("ImportedTransfers properly marks number of accounts imported")
	{
		require.Equal(t, int64(0), entity.importedTransfers.Count())
		entity.ImportedTransfers(8)
		assert.Equal(t, int64(8), entity.importedTransfers.Count())
	}

	t.Log("ExportedTransfers properly marks number of accounts exported")
	{
		require.Equal(t, int64(0), entity.exportedTransfers.Count())
		entity.ExportedTransfers(16)
		assert.Equal(t, int64(16), entity.exportedTransfers.Count())
	}
}
