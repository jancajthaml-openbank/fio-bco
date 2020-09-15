package metrics

import (
	"encoding/json"
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
			createTokenLatency: metrics.NewTimer(),
			deleteTokenLatency: metrics.NewTimer(),
			getTokenLatency:    metrics.NewTimer(),
		}

		entity.getTokenLatency.Update(time.Duration(1))
		entity.createTokenLatency.Update(time.Duration(2))
		entity.deleteTokenLatency.Update(time.Duration(3))

		actual, err := entity.MarshalJSON()
		require.Nil(t, err)

		aux := &struct {
			GetTokenLatency    float64 `json:"getTokenLatency"`
			CreateTokenLatency float64 `json:"createTokenLatency"`
			DeleteTokenLatency float64 `json:"deleteTokenLatency"`
		}{}

		require.Nil(t, json.Unmarshal(actual, &aux))

		assert.Equal(t, float64(1), aux.GetTokenLatency)
		assert.Equal(t, float64(2), aux.CreateTokenLatency)
		assert.Equal(t, float64(3), aux.DeleteTokenLatency)
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
			createTokenLatency: metrics.NewTimer(),
			deleteTokenLatency: metrics.NewTimer(),
			getTokenLatency:    metrics.NewTimer(),
		}

		data := []byte("{")
		assert.NotNil(t, entity.UnmarshalJSON(data))
	}

	t.Log("happy path")
	{
		entity := Metrics{
			createTokenLatency: metrics.NewTimer(),
			deleteTokenLatency: metrics.NewTimer(),
			getTokenLatency:    metrics.NewTimer(),
		}

		data := []byte("{\"getTokenLatency\":1,\"createTokenLatency\":2,\"deleteTokenLatency\":3}")
		require.Nil(t, entity.UnmarshalJSON(data))

		assert.Equal(t, float64(1), entity.getTokenLatency.Percentile(0.95))
		assert.Equal(t, float64(2), entity.createTokenLatency.Percentile(0.95))
		assert.Equal(t, float64(3), entity.deleteTokenLatency.Percentile(0.95))
	}
}

func TestPersist(t *testing.T) {

	t.Log("error when caller is nil")
	{
		var entity *Metrics
		assert.NotNil(t, entity.Persist())
	}

	t.Log("error when marshaling fails")
	{
		entity := Metrics{}
		assert.NotNil(t, entity.Persist())
	}

	t.Log("happy path")
	{
		defer os.Remove("/tmp/metrics.json")

		entity := Metrics{
			storage:            localfs.NewPlaintextStorage("/tmp"),
			createTokenLatency: metrics.NewTimer(),
			deleteTokenLatency: metrics.NewTimer(),
			getTokenLatency:    metrics.NewTimer(),
		}

		require.Nil(t, entity.Persist())

		expected, err := entity.MarshalJSON()
		require.Nil(t, err)

		actual, err := ioutil.ReadFile("/tmp/metrics.json")
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
		defer os.Remove("/tmp/metrics.json")

		old := Metrics{
			createTokenLatency: metrics.NewTimer(),
			deleteTokenLatency: metrics.NewTimer(),
			getTokenLatency:    metrics.NewTimer(),
		}

		old.getTokenLatency.Update(time.Duration(1))
		old.createTokenLatency.Update(time.Duration(2))
		old.deleteTokenLatency.Update(time.Duration(3))

		data, err := old.MarshalJSON()
		require.Nil(t, err)

		require.Nil(t, ioutil.WriteFile("/tmp/metrics.json", data, 0444))

		entity := Metrics{
			storage:            localfs.NewPlaintextStorage("/tmp"),
			createTokenLatency: metrics.NewTimer(),
			deleteTokenLatency: metrics.NewTimer(),
			getTokenLatency:    metrics.NewTimer(),
		}

		require.Nil(t, entity.Hydrate())

		assert.Equal(t, float64(1), entity.getTokenLatency.Percentile(0.95))
		assert.Equal(t, float64(2), entity.createTokenLatency.Percentile(0.95))
		assert.Equal(t, float64(3), entity.deleteTokenLatency.Percentile(0.95))
	}
}
