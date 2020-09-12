// Copyright (c) 2016-2020, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"bytes"
	"fmt"
	"github.com/jancajthaml-openbank/fio-bco-import/utils"
	"os"
	"strconv"
	"time"
)

// MarshalJSON serializes Metrics as json bytes
func (metrics *Metrics) MarshalJSON() ([]byte, error) {
	if metrics == nil {
		return nil, fmt.Errorf("cannot marshall nil")
	}

	if metrics.createdTokens == nil || metrics.deletedTokens == nil ||
		metrics.syncLatency == nil || metrics.importedTransfers == nil ||
		metrics.importedTransactions == nil {
		return nil, fmt.Errorf("cannot marshall nil references")
	}

	var buffer bytes.Buffer

	buffer.WriteString("{\"createdTokens\":")
	buffer.WriteString(strconv.FormatInt(metrics.createdTokens.Count(), 10))
	buffer.WriteString(",\"deletedTokens\":")
	buffer.WriteString(strconv.FormatInt(metrics.deletedTokens.Count(), 10))
	buffer.WriteString(",\"syncLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.syncLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importedTransfers\":")
	buffer.WriteString(strconv.FormatInt(metrics.importedTransfers.Count(), 10))
	buffer.WriteString(",\"importedTransactions\":")
	buffer.WriteString(strconv.FormatInt(metrics.importedTransactions.Count(), 10))
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON deserializes Metrics from json bytes
func (metrics *Metrics) UnmarshalJSON(data []byte) error {
	if metrics == nil {
		return fmt.Errorf("cannot unmarshall to nil")
	}

	if metrics.createdTokens == nil || metrics.deletedTokens == nil ||
		metrics.syncLatency == nil || metrics.importedTransfers == nil ||
		metrics.importedTransactions == nil {
		return fmt.Errorf("cannot unmarshall to nil references")
	}

	aux := &struct {
		CreatedTokens        int64   `json:"createdTokens"`
		DeletedTokens        int64   `json:"deletedTokens"`
		SyncLatency          float64 `json:"syncLatency"`
		ImportedTransfers    int64   `json:"importedTransfers"`
		ImportedTransactions int64   `json:"importedTransactions"`
	}{}

	if err := utils.JSON.Unmarshal(data, &aux); err != nil {
		return err
	}

	metrics.createdTokens.Clear()
	metrics.createdTokens.Inc(aux.CreatedTokens)
	metrics.deletedTokens.Clear()
	metrics.deletedTokens.Inc(aux.DeletedTokens)
	metrics.syncLatency.Update(time.Duration(aux.SyncLatency))
	metrics.importedTransfers.Mark(aux.ImportedTransfers)
	metrics.importedTransactions.Mark(aux.ImportedTransactions)

	return nil
}

// Persist saved metrics state to storage
func (metrics *Metrics) Persist() error {
	if metrics == nil {
		return fmt.Errorf("cannot persist nil reference")
	}
	data, err := utils.JSON.Marshal(metrics)
	if err != nil {
		log.Warn().Msgf("unable to marshall metrics %+v", err)
		return err
	}
	err = metrics.storage.WriteFile("metrics."+metrics.tenant+".json", data)
	if err != nil {
		log.Warn().Msgf("unable to persist metrics %+v", err)
		return err
	}
	err = os.Chmod(metrics.storage.Root+"/metrics."+metrics.tenant+".json", 0644)
	if err != nil {
		return err
	}
	return nil
}

// Hydrate loads metrics state from storage
func (metrics *Metrics) Hydrate() error {
	if metrics == nil {
		return fmt.Errorf("cannot hydrate nil reference")
	}
	data, err := metrics.storage.ReadFileFully("metrics." + metrics.tenant + ".json")
	if err != nil {
		return err
	}
	err = utils.JSON.Unmarshal(data, metrics)
	if err != nil {
		return err
	}
	return nil
}
