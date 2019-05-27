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
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/utils"
)

// MarshalJSON serialises Metrics as json preserving uint64
func (entity *Metrics) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString("{\"createdTokens\":")
	buffer.WriteString(strconv.FormatInt(entity.createdTokens.Count(), 10))
	buffer.WriteString(",\"deletedTokens\":")
	buffer.WriteString(strconv.FormatInt(entity.deletedTokens.Count(), 10))
	buffer.WriteString(",\"syncLatency\":")
	buffer.WriteString(strconv.FormatFloat(entity.syncLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importAccountLatency\":")
	buffer.WriteString(strconv.FormatFloat(entity.importAccountLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"exportAccountLatency\":")
	buffer.WriteString(strconv.FormatFloat(entity.exportAccountLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importTransactionLatency\":")
	buffer.WriteString(strconv.FormatFloat(entity.importTransactionLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"exportTransactionLatency\":")
	buffer.WriteString(strconv.FormatFloat(entity.exportTransactionLatency.Percentile(0.95), 'f', -1, 64))
	buffer.WriteString(",\"importedAccounts\":")
	buffer.WriteString(strconv.FormatInt(entity.importedAccounts.Count(), 10))
	buffer.WriteString(",\"exportedAccounts\":")
	buffer.WriteString(strconv.FormatInt(entity.exportedAccounts.Count(), 10))
	buffer.WriteString(",\"importedTransfers\":")
	buffer.WriteString(strconv.FormatInt(entity.importedTransfers.Count(), 10))
	buffer.WriteString(",\"exportedTransfers\":")
	buffer.WriteString(strconv.FormatInt(entity.exportedTransfers.Count(), 10))
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshal json of Metrics entity
func (entity *Metrics) UnmarshalJSON(data []byte) error {
	if entity == nil {
		return fmt.Errorf("cannot unmarshall to nil pointer")
	}
	all := struct {
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
	err := utils.JSON.Unmarshal(data, &all)
	if err != nil {
		return err
	}

	entity.createdTokens.Clear()
	entity.createdTokens.Inc(all.CreatedTokens)

	entity.deletedTokens.Clear()
	entity.deletedTokens.Inc(all.DeletedTokens)

	entity.syncLatency.Update(time.Duration(all.SyncLatency))
	entity.importAccountLatency.Update(time.Duration(all.ImportAccountLatency))
	entity.exportAccountLatency.Update(time.Duration(all.ExportAccountLatency))
	entity.importTransactionLatency.Update(time.Duration(all.ImportTransactionLatency))
	entity.exportTransactionLatency.Update(time.Duration(all.ExportTransactionLatency))

	entity.importedAccounts.Mark(all.ImportedAccounts)
	entity.exportedAccounts.Mark(all.ExportedAccounts)
	entity.importedTransfers.Mark(all.ImportedTransfers)
	entity.exportedTransfers.Mark(all.ExportedTransfers)

	return nil
}

func (metrics *Metrics) Persist() error {
	if metrics == nil {
		return fmt.Errorf("cannot persist nil reference")
	}
	tempFile := metrics.output + "_temp"
	data, err := utils.JSON.Marshal(metrics)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := os.Rename(tempFile, metrics.output); err != nil {
		return err
	}
	return nil
}

func (metrics *Metrics) Hydrate() error {
	if metrics == nil {
		return fmt.Errorf("cannot hydrate nil reference")
	}
	f, err := os.OpenFile(metrics.output, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	buf := make([]byte, fi.Size())
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	err = utils.JSON.Unmarshal(buf, metrics)
	if err != nil {
		return err
	}
	return nil
}
