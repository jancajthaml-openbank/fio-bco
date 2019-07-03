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
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// TokenCreated increments token created by one
func (metrics *Metrics) TokenCreated() {
	metrics.createdTokens.Inc(1)
}

// TokenDeleted increments token deleted by one
func (metrics *Metrics) TokenDeleted() {
	metrics.deletedTokens.Inc(1)
}

func (metrics *Metrics) TimeSyncLatency(f func()) {
	metrics.syncLatency.Time(f)
}

func (metrics *Metrics) TimeImportAccount(f func()) {
	metrics.importAccountLatency.Time(f)
}

func (metrics *Metrics) TimeExportAccount(f func()) {
	metrics.exportAccountLatency.Time(f)
}

func (metrics *Metrics) TimeImportTransaction(f func()) {
	metrics.importTransactionLatency.Time(f)
}

func (metrics *Metrics) TimeExportTransaction(f func()) {
	metrics.exportTransactionLatency.Time(f)
}

func (metrics *Metrics) ImportedAccounts(num int64) {
	metrics.importedAccounts.Mark(num)
}

func (metrics *Metrics) ExportedAccounts(num int64) {
	metrics.exportedAccounts.Mark(num)
}

func (metrics *Metrics) ImportedTransfers(num int64) {
	metrics.importedTransfers.Mark(num)
}

func (metrics *Metrics) ExportedTransfers(num int64) {
	metrics.exportedTransfers.Mark(num)
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

	ticker := time.NewTicker(metrics.refreshRate)
	defer ticker.Stop()

	if err := metrics.Hydrate(); err != nil {
		log.Warn(err.Error())
	}
	metrics.MarkReady()

	select {
	case <-metrics.CanStart:
		break
	case <-metrics.Done():
		return
	}

	log.Infof("Start metrics daemon, update each %v into %v", metrics.refreshRate, metrics.output)

	for {
		select {
		case <-metrics.Done():
			log.Info("Stopping metrics daemon")
			metrics.Persist()
			log.Info("Stop metrics daemon")
			return
		case <-ticker.C:
			metrics.Persist()
		}
	}
}
