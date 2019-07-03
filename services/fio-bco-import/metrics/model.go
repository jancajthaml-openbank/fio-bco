// Copyright (c) 2016-2019, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/utils"
	metrics "github.com/rcrowley/go-metrics"
)

// Metrics represents metrics subroutine
type Metrics struct {
	utils.DaemonSupport
	output                   string
	createdTokens            metrics.Counter
	deletedTokens            metrics.Counter
	refreshRate              time.Duration
	syncLatency              metrics.Timer
	importAccountLatency     metrics.Timer
	exportAccountLatency     metrics.Timer
	importTransactionLatency metrics.Timer
	exportTransactionLatency metrics.Timer
	importedAccounts         metrics.Meter
	exportedAccounts         metrics.Meter
	importedTransfers        metrics.Meter
	exportedTransfers        metrics.Meter
}

// NewMetrics returns metrics fascade
func NewMetrics(ctx context.Context, output string, refreshRate time.Duration) Metrics {
	return Metrics{
		DaemonSupport:            utils.NewDaemonSupport(ctx),
		output:                   output,
		refreshRate:              refreshRate,
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
}

// MarshalJSON serialises Metrics as json bytes
func (metrics *Metrics) MarshalJSON() ([]byte, error) {
	if metrics == nil {
		return nil, fmt.Errorf("cannot marshall nil")
	}

	if metrics.createdTokens == nil || metrics.deletedTokens == nil ||
		metrics.syncLatency == nil || metrics.importAccountLatency == nil ||
		metrics.exportAccountLatency == nil || metrics.importTransactionLatency == nil ||
		metrics.exportTransactionLatency == nil || metrics.importedAccounts == nil ||
		metrics.exportedAccounts == nil || metrics.importedTransfers == nil ||
		metrics.exportedTransfers == nil {
		return nil, fmt.Errorf("cannot marshall nil references")
	}

	var buffer bytes.Buffer

	buffer.WriteString("{\"createdTokens\":")
	buffer.WriteString(strconv.FormatInt(metrics.createdTokens.Count(), 10))
	buffer.WriteString(",\"deletedTokens\":")
	buffer.WriteString(strconv.FormatInt(metrics.deletedTokens.Count(), 10))
	buffer.WriteString(",\"syncLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.syncLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importAccountLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.importAccountLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"exportAccountLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.exportAccountLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importTransactionLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.importTransactionLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"exportTransactionLatency\":")
	buffer.WriteString(strconv.FormatFloat(metrics.exportTransactionLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importedAccounts\":")
	buffer.WriteString(strconv.FormatInt(metrics.importedAccounts.Count(), 10))
	buffer.WriteString(",\"exportedAccounts\":")
	buffer.WriteString(strconv.FormatInt(metrics.exportedAccounts.Count(), 10))
	buffer.WriteString(",\"importedTransfers\":")
	buffer.WriteString(strconv.FormatInt(metrics.importedTransfers.Count(), 10))
	buffer.WriteString(",\"exportedTransfers\":")
	buffer.WriteString(strconv.FormatInt(metrics.exportedTransfers.Count(), 10))
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON deserializes Metrics from json bytes
func (metrics *Metrics) UnmarshalJSON(data []byte) error {
	if metrics == nil {
		return fmt.Errorf("cannot unmarshall to nil")
	}

	if metrics.createdTokens == nil || metrics.deletedTokens == nil ||
		metrics.syncLatency == nil || metrics.importAccountLatency == nil ||
		metrics.exportAccountLatency == nil || metrics.importTransactionLatency == nil ||
		metrics.exportTransactionLatency == nil || metrics.importedAccounts == nil ||
		metrics.exportedAccounts == nil || metrics.importedTransfers == nil ||
		metrics.exportedTransfers == nil {
		return fmt.Errorf("cannot unmarshall to nil references")
	}

	aux := &struct {
		CreatedTokens            int64   `json:"createdTokens"`
		DeletedTokens            int64   `json:"deletedTokens"`
		SyncLatency              float64 `json:"syncLatency"`
		ImportAccountLatency     float64 `json:"importAccountLatency"`
		ExportAccountLatency     float64 `json:"exportAccountLatency"`
		ImportTransactionLatency float64 `json:"importTransactionLatency"`
		ExportTransactionLatency float64 `json:"exportTransactionLatency"`
		ImportedAccounts         int64   `json:"importedAccounts"`
		ExportedAccounts         int64   `json:"exportedAccounts"`
		ImportedTransfers        int64   `json:"importedTransfers"`
		ExportedTransfers        int64   `json:"exportedTransfers"`
	}{}

	if err := utils.JSON.Unmarshal(data, &aux); err != nil {
		return err
	}

	metrics.createdTokens.Clear()
	metrics.createdTokens.Inc(aux.CreatedTokens)
	metrics.deletedTokens.Clear()
	metrics.deletedTokens.Inc(aux.DeletedTokens)
	metrics.syncLatency.Update(time.Duration(aux.SyncLatency))
	metrics.importAccountLatency.Update(time.Duration(aux.ImportAccountLatency))
	metrics.exportAccountLatency.Update(time.Duration(aux.ExportAccountLatency))
	metrics.importTransactionLatency.Update(time.Duration(aux.ImportTransactionLatency))
	metrics.exportTransactionLatency.Update(time.Duration(aux.ExportTransactionLatency))
	metrics.importedAccounts.Mark(aux.ImportedAccounts)
	metrics.exportedAccounts.Mark(aux.ExportedAccounts)
	metrics.importedTransfers.Mark(aux.ImportedTransfers)
	metrics.exportedTransfers.Mark(aux.ExportedTransfers)

	return nil
}
