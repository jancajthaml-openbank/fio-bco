package metrics

import (
	"testing"
	"time"

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
			createTokenLatency: metrics.NewTimer(),
			deleteTokenLatency: metrics.NewTimer(),
			getTokenLatency:    metrics.NewTimer(),
		}

		entity.getTokenLatency.Update(time.Duration(1))
		entity.createTokenLatency.Update(time.Duration(2))
		entity.deleteTokenLatency.Update(time.Duration(3))

		actual, err := entity.MarshalJSON()

		require.Nil(t, err)

		data := []byte("{\"getTokenLatency\":1,\"createTokenLatency\":2,\"deleteTokenLatency\":3}")

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
