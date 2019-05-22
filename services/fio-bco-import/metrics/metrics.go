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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/utils"

	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

// Metrics represents metrics subroutine
type Metrics struct {
	utils.DaemonSupport
	output                   string
	tenant                   string
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

// Snapshot holds metrics snapshot status
type Snapshot struct {
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
}

// NewMetrics returns metrics fascade
func NewMetrics(ctx context.Context, tenant string, output string, refreshRate time.Duration) Metrics {
	return Metrics{
		DaemonSupport:            utils.NewDaemonSupport(ctx),
		output:                   output,
		tenant:                   tenant,
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

// NewSnapshot returns metrics snapshot
func NewSnapshot(metrics Metrics) Snapshot {
	return Snapshot{
		CreatedTokens:            metrics.createdTokens.Count(),
		DeletedTokens:            metrics.deletedTokens.Count(),
		SyncLatency:              metrics.syncLatency.Percentile(0.95),
		ImportAccountLatency:     metrics.importAccountLatency.Percentile(0.95),
		ExportAccountLatency:     metrics.exportAccountLatency.Percentile(0.95),
		ImportTransactionLatency: metrics.importTransactionLatency.Percentile(0.95),
		ExportTransactionLatency: metrics.exportTransactionLatency.Percentile(0.95),
		ImportedAccounts:         metrics.importedAccounts.Count(),
		ExportedAccounts:         metrics.exportedAccounts.Count(),
		ImportedTransfers:        metrics.importedTransfers.Count(),
		ExportedTransfers:        metrics.exportedTransfers.Count(),
	}
}

// TokenCreated increments token created by one
func (metrics Metrics) TokenCreated() {
	metrics.createdTokens.Inc(1)
}

// TokenDeleted increments token deleted by one
func (metrics Metrics) TokenDeleted() {
	metrics.deletedTokens.Inc(1)
}

func (metrics Metrics) TimeSyncLatency(f func()) {
	metrics.syncLatency.Time(f)
}

func (metrics Metrics) TimeImportAccount(f func()) {
	metrics.importAccountLatency.Time(f)
}

func (metrics Metrics) TimeExportAccount(f func()) {
	metrics.exportAccountLatency.Time(f)
}

func (metrics Metrics) TimeImportTransaction(f func()) {
	metrics.importTransactionLatency.Time(f)
}

func (metrics Metrics) TimeExportTransaction(f func()) {
	metrics.exportTransactionLatency.Time(f)
}

func (metrics Metrics) ImportedAccounts(num int64) {
	metrics.importedAccounts.Mark(num)
}

func (metrics Metrics) ExportedAccounts(num int64) {
	metrics.exportedAccounts.Mark(num)
}

func (metrics Metrics) ImportedTransfers(num int64) {
	metrics.importedTransfers.Mark(num)
}

func (metrics Metrics) ExportedTransfers(num int64) {
	metrics.exportedTransfers.Mark(num)
}

func (metrics Metrics) persist(filename string) {
	tempFile := filename + "_temp"

	data, err := utils.JSON.Marshal(NewSnapshot(metrics))
	if err != nil {
		log.Warnf("unable to create serialize metrics with error: %v", err)
		return
	}
	f, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Warnf("unable to create file with error: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		log.Warnf("unable to write file with error: %v", err)
		return
	}

	if err := os.Rename(tempFile, filename); err != nil {
		log.Warnf("unable to move file with error: %v", err)
		return
	}

	return
}

func getFilename(path, tenant string) string {
	if tenant == "" {
		return path
	}

	dirname := filepath.Dir(path)
	ext := filepath.Ext(path)
	filename := filepath.Base(path)
	filename = filename[:len(filename)-len(ext)]

	return dirname + "/" + filename + ".import." + tenant + ext
}

// WaitReady wait for metrics to be ready
func (metrics Metrics) WaitReady(deadline time.Duration) (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch x := e.(type) {
			case string:
				err = fmt.Errorf(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic")
			}
		}
	}()

	ticker := time.NewTicker(deadline)
	select {
	case <-metrics.IsReady:
		ticker.Stop()
		err = nil
		return
	case <-ticker.C:
		err = fmt.Errorf("daemon was not ready within %v seconds", deadline)
		return
	}
}

// Start handles everything needed to start metrics daemon
func (metrics Metrics) Start() {
	defer metrics.MarkDone()

	if metrics.output == "" {
		log.Warnf("no metrics output defined, skipping metrics persistence")
		metrics.MarkReady()
		return
	}

	output := getFilename(metrics.output, metrics.tenant)
	ticker := time.NewTicker(metrics.refreshRate)
	defer ticker.Stop()

	metrics.MarkReady()

	select {
	case <-metrics.CanStart:
		break
	case <-metrics.Done():
		return
	}

	log.Infof("Start metrics daemon, update each %v into %v", metrics.refreshRate, output)

	for {
		select {
		case <-metrics.Done():
			log.Info("Stopping metrics daemon")
			metrics.persist(output)
			log.Info("Stop metrics daemon")
			return
		case <-ticker.C:
			metrics.persist(output)
		}
	}
}
