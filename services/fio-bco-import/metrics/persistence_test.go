package metrics

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersist(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.EqualError(t, entity.Persist(), "cannot persist nil reference")
	}

	t.Log("error when marshalling fails")
	{
		entity := Metrics{}
		assert.EqualError(t, entity.Persist(), "cannot marshall nil references")
	}

	t.Log("error when cannot open tempfile for writing")
	{
		entity := Metrics{
			output:                   "/sys/kernel/security",
			createdTokens:            metrics.NewCounter(),
			deletedTokens:            metrics.NewCounter(),
			syncLatency:              metrics.NewTimer(),
			importAccountLatency:     metrics.NewTimer(),
			exportAccountLatency:     metrics.NewTimer(),
			importTransactionLatency: metrics.NewTimer(),
			exportTransactionLatency: metrics.NewTimer(),
			importedAccounts:         metrics.NewMeter(),
			exportedAccounts:         metrics.NewMeter(),
			importedTransfers:        metrics.NewMeter(),
			exportedTransfers:        metrics.NewMeter(),
		}

		assert.NotNil(t, entity.Persist())
	}

	t.Log("happy path")
	{
		tmpfile, err := ioutil.TempFile(os.TempDir(), "test_metrics_persist")

		require.Nil(t, err)
		defer os.Remove(tmpfile.Name())

		entity := Metrics{
			output:                   tmpfile.Name(),
			createdTokens:            metrics.NewCounter(),
			deletedTokens:            metrics.NewCounter(),
			syncLatency:              metrics.NewTimer(),
			importAccountLatency:     metrics.NewTimer(),
			exportAccountLatency:     metrics.NewTimer(),
			importTransactionLatency: metrics.NewTimer(),
			exportTransactionLatency: metrics.NewTimer(),
			importedAccounts:         metrics.NewMeter(),
			exportedAccounts:         metrics.NewMeter(),
			importedTransfers:        metrics.NewMeter(),
			exportedTransfers:        metrics.NewMeter(),
		}

		require.Nil(t, entity.Persist())

		expected, err := entity.MarshalJSON()
		require.Nil(t, err)

		actual, err := ioutil.ReadFile(tmpfile.Name())
		require.Nil(t, err)

		assert.Equal(t, expected, actual)
	}
}

func TestHydrate(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.EqualError(t, entity.Hydrate(), "cannot hydrate nil reference")
	}

	t.Log("happy path")
	{
		tmpfile, err := ioutil.TempFile(os.TempDir(), "test_metrics_hydrate")

		require.Nil(t, err)
		defer os.Remove(tmpfile.Name())

		old := Metrics{
			createdTokens:            metrics.NewCounter(),
			deletedTokens:            metrics.NewCounter(),
			syncLatency:              metrics.NewTimer(),
			importAccountLatency:     metrics.NewTimer(),
			exportAccountLatency:     metrics.NewTimer(),
			importTransactionLatency: metrics.NewTimer(),
			exportTransactionLatency: metrics.NewTimer(),
			importedAccounts:         metrics.NewMeter(),
			exportedAccounts:         metrics.NewMeter(),
			importedTransfers:        metrics.NewMeter(),
			exportedTransfers:        metrics.NewMeter(),
		}

		old.createdTokens.Inc(1)
		old.deletedTokens.Inc(2)
		old.syncLatency.Update(time.Duration(3))
		old.importAccountLatency.Update(time.Duration(4))
		old.exportAccountLatency.Update(time.Duration(5))
		old.importTransactionLatency.Update(time.Duration(6))
		old.exportTransactionLatency.Update(time.Duration(7))
		old.importedAccounts.Mark(8)
		old.exportedAccounts.Mark(9)
		old.importedTransfers.Mark(10)
		old.exportedTransfers.Mark(11)

		data, err := old.MarshalJSON()
		require.Nil(t, err)

		require.Nil(t, ioutil.WriteFile(tmpfile.Name(), data, 0444))

		entity := Metrics{
			output:                   tmpfile.Name(),
			createdTokens:            metrics.NewCounter(),
			deletedTokens:            metrics.NewCounter(),
			syncLatency:              metrics.NewTimer(),
			importAccountLatency:     metrics.NewTimer(),
			exportAccountLatency:     metrics.NewTimer(),
			importTransactionLatency: metrics.NewTimer(),
			exportTransactionLatency: metrics.NewTimer(),
			importedAccounts:         metrics.NewMeter(),
			exportedAccounts:         metrics.NewMeter(),
			importedTransfers:        metrics.NewMeter(),
			exportedTransfers:        metrics.NewMeter(),
		}

		require.Nil(t, entity.Hydrate())

		assert.Equal(t, int64(1), entity.createdTokens.Count())
		assert.Equal(t, int64(2), entity.deletedTokens.Count())
		assert.Equal(t, float64(3), entity.syncLatency.Percentile(0.95))
		assert.Equal(t, float64(4), entity.importAccountLatency.Percentile(0.95))
		assert.Equal(t, float64(5), entity.exportAccountLatency.Percentile(0.95))
		assert.Equal(t, float64(6), entity.importTransactionLatency.Percentile(0.95))
		assert.Equal(t, float64(7), entity.exportTransactionLatency.Percentile(0.95))
		assert.Equal(t, int64(8), entity.importedAccounts.Count())
		assert.Equal(t, int64(9), entity.exportedAccounts.Count())
		assert.Equal(t, int64(10), entity.importedTransfers.Count())
		assert.Equal(t, int64(11), entity.exportedTransfers.Count())
	}
}
