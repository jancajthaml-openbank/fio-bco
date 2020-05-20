package metrics

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	localfs "github.com/jancajthaml-openbank/local-fs"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		_, err := entity.MarshalJSON()
		assert.EqualError(t, err, "cannot marshall nil")
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		_, err := entity.MarshalJSON()
		assert.EqualError(t, err, "cannot marshall nil references")
	}

	t.Log("happy path")
	{
		entity := Metrics{
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

		entity.createdTokens.Inc(1)
		entity.deletedTokens.Inc(2)
		entity.syncLatency.Update(time.Duration(3))
		entity.importAccountLatency.Update(time.Duration(4))
		entity.exportAccountLatency.Update(time.Duration(5))
		entity.importTransactionLatency.Update(time.Duration(6))
		entity.exportTransactionLatency.Update(time.Duration(7))
		entity.importedAccounts.Mark(8)
		entity.exportedAccounts.Mark(9)
		entity.importedTransfers.Mark(10)
		entity.exportedTransfers.Mark(11)

		actual, err := entity.MarshalJSON()

		require.Nil(t, err)

		data := []byte("{\"createdTokens\":1,\"deletedTokens\":2,\"syncLatency\":3,\"importAccountLatency\":4,\"exportAccountLatency\":5,\"importTransactionLatency\":6,\"exportTransactionLatency\":7,\"importedAccounts\":8,\"exportedAccounts\":9,\"importedTransfers\":10,\"exportedTransfers\":11}")

		assert.Equal(t, data, actual)
	}
}

func TestUnmarshalJSON(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		err := entity.UnmarshalJSON([]byte(""))
		assert.EqualError(t, err, "cannot unmarshall to nil")
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		err := entity.UnmarshalJSON([]byte(""))
		assert.EqualError(t, err, "cannot unmarshall to nil references")
	}

	t.Log("error on malformed data")
	{
		entity := Metrics{
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

		data := []byte("{")
		assert.NotNil(t, entity.UnmarshalJSON(data))
	}

	t.Log("happy path")
	{
		entity := Metrics{
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

		data := []byte("{\"createdTokens\":1,\"deletedTokens\":2,\"syncLatency\":3,\"importAccountLatency\":4,\"exportAccountLatency\":5,\"importTransactionLatency\":6,\"exportTransactionLatency\":7,\"importedAccounts\":8,\"exportedAccounts\":9,\"importedTransfers\":10,\"exportedTransfers\":11}")
		require.Nil(t, entity.UnmarshalJSON(data))

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


	t.Log("happy path")
	{
		defer os.Remove("/tmp/metrics.json")

		entity := Metrics{
			storage:                  localfs.NewPlaintextStorage("/tmp"),
			tenant:                   "1",
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

		actual, err := ioutil.ReadFile("/tmp/metrics.1.json")
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
		defer os.Remove("/tmp/metrics.1.json")

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

		require.Nil(t, ioutil.WriteFile("/tmp/metrics.1.json", data, 0444))

		entity := Metrics{
			storage:                  localfs.NewPlaintextStorage("/tmp"),
			tenant:                   "1",
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
