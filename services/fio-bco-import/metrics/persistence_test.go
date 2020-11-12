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
		assert.NotNil(t, err)
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		_, err := entity.MarshalJSON()
		assert.NotNil(t, err)
	}

	t.Log("happy path")
	{
		entity := Metrics{
			createdTokens:        metrics.NewCounter(),
			deletedTokens:        metrics.NewCounter(),
			syncLatency:          metrics.NewTimer(),
			importedTransfers:    metrics.NewMeter(),
			importedTransactions: metrics.NewMeter(),
		}

		entity.createdTokens.Inc(1)
		entity.deletedTokens.Inc(2)
		entity.syncLatency.Update(time.Duration(3))
		entity.importedTransfers.Mark(4)
		entity.importedTransactions.Mark(5)

		actual, err := entity.MarshalJSON()

		require.Nil(t, err)

		data := []byte("{\"createdTokens\":1,\"deletedTokens\":2,\"syncLatency\":3,\"importedTransfers\":4,\"importedTransactions\":5}")

		assert.Equal(t, data, actual)
	}
}

func TestUnmarshalJSON(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		err := entity.UnmarshalJSON([]byte(""))
		assert.NotNil(t, err)
	}

	t.Log("error when values are nil")
	{
		entity := Metrics{}
		err := entity.UnmarshalJSON([]byte(""))
		assert.NotNil(t, err)
	}

	t.Log("error on malformed data")
	{
		entity := Metrics{
			createdTokens:        metrics.NewCounter(),
			deletedTokens:        metrics.NewCounter(),
			syncLatency:          metrics.NewTimer(),
			importedTransfers:    metrics.NewMeter(),
			importedTransactions: metrics.NewMeter(),
		}

		data := []byte("{")
		assert.NotNil(t, entity.UnmarshalJSON(data))
	}

	t.Log("happy path")
	{
		entity := Metrics{
			createdTokens:        metrics.NewCounter(),
			deletedTokens:        metrics.NewCounter(),
			syncLatency:          metrics.NewTimer(),
			importedTransfers:    metrics.NewMeter(),
			importedTransactions: metrics.NewMeter(),
		}

		data := []byte("{\"createdTokens\":1,\"deletedTokens\":2,\"syncLatency\":3,\"importedTransfers\":4,\"importedTransactions\":5}")
		require.Nil(t, entity.UnmarshalJSON(data))

		assert.Equal(t, int64(1), entity.createdTokens.Count())
		assert.Equal(t, int64(2), entity.deletedTokens.Count())
		assert.Equal(t, float64(3), entity.syncLatency.Percentile(0.95))
		assert.Equal(t, int64(4), entity.importedTransfers.Count())
		assert.Equal(t, int64(5), entity.importedTransactions.Count())
	}
}

func TestPersist(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.EqualError(t, entity.Persist(), "cannot persist nil reference")
	}

	t.Log("error when marshaling fails")
	{
		entity := Metrics{}
		assert.NotNil(t, entity.Persist())
	}

	t.Log("happy path")
	{
		defer os.Remove("/tmp/metrics.json")

		storage, _ := localfs.NewPlaintextStorage("/tmp")

		entity := Metrics{
			storage:              storage,
			tenant:               "1",
			createdTokens:        metrics.NewCounter(),
			deletedTokens:        metrics.NewCounter(),
			syncLatency:          metrics.NewTimer(),
			importedTransfers:    metrics.NewMeter(),
			importedTransactions: metrics.NewMeter(),
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
			createdTokens:        metrics.NewCounter(),
			deletedTokens:        metrics.NewCounter(),
			syncLatency:          metrics.NewTimer(),
			importedTransfers:    metrics.NewMeter(),
			importedTransactions: metrics.NewMeter(),
		}

		old.createdTokens.Inc(1)
		old.deletedTokens.Inc(2)
		old.syncLatency.Update(time.Duration(3))
		old.importedTransfers.Mark(4)
		old.importedTransactions.Mark(5)

		data, err := old.MarshalJSON()
		require.Nil(t, err)

		require.Nil(t, ioutil.WriteFile("/tmp/metrics.1.json", data, 0444))

		storage, _ := localfs.NewPlaintextStorage("/tmp")

		entity := Metrics{
			storage:              storage,
			tenant:               "1",
			createdTokens:        metrics.NewCounter(),
			deletedTokens:        metrics.NewCounter(),
			syncLatency:          metrics.NewTimer(),
			importedTransfers:    metrics.NewMeter(),
			importedTransactions: metrics.NewMeter(),
		}

		require.Nil(t, entity.Hydrate())

		assert.Equal(t, int64(1), entity.createdTokens.Count())
		assert.Equal(t, int64(2), entity.deletedTokens.Count())
		assert.Equal(t, float64(3), entity.syncLatency.Percentile(0.95))
		assert.Equal(t, int64(4), entity.importedTransfers.Count())
		assert.Equal(t, int64(5), entity.importedTransactions.Count())
	}
}
