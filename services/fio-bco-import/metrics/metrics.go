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
	"context"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/utils"
	localfs "github.com/jancajthaml-openbank/local-fs"
	metrics "github.com/rcrowley/go-metrics"
)

// Metrics holds metrics counters
type Metrics struct {
	utils.DaemonSupport
	storage              localfs.PlaintextStorage
	tenant               string
	refreshRate          time.Duration
	createdTokens        metrics.Counter
	deletedTokens        metrics.Counter
	syncLatency          metrics.Timer
	importedTransfers    metrics.Meter
	importedTransactions metrics.Meter
}

// NewMetrics returns blank metrics holder
func NewMetrics(ctx context.Context, output string, tenant string, refreshRate time.Duration) Metrics {
	return Metrics{
		DaemonSupport:        utils.NewDaemonSupport(ctx, "metrics"),
		storage:              localfs.NewPlaintextStorage(output),
		tenant:               tenant,
		refreshRate:          refreshRate,
		createdTokens:        metrics.NewCounter(),
		deletedTokens:        metrics.NewCounter(),
		syncLatency:          metrics.NewTimer(),
		importedTransfers:    metrics.NewMeter(),
		importedTransactions: metrics.NewMeter(),
	}
}

// TokenCreated increments token created by one
func (metrics *Metrics) TokenCreated() {
	metrics.createdTokens.Inc(1)
}

// TokenDeleted increments token deleted by one
func (metrics *Metrics) TokenDeleted() {
	metrics.deletedTokens.Inc(1)
}

// TimeSyncLatency measures time of sync duration
func (metrics *Metrics) TimeSyncLatency(f func()) {
	metrics.syncLatency.Time(f)
}

// TransactionImported increments transactions importer by one
func (metrics *Metrics) TransactionImported() {
	metrics.importedTransactions.Mark(1)
}

// TransfersImported increments transfers importer by given num
func (metrics *Metrics) TransfersImported(num int64) {
	metrics.importedTransfers.Mark(num)
}

// Start handles everything needed to start metrics daemon
func (metrics *Metrics) Start() {
	ticker := time.NewTicker(metrics.refreshRate)
	defer ticker.Stop()

	if err := metrics.Hydrate(); err != nil {
		log.Warn().Msg(err.Error())
	}

	metrics.Persist()
	metrics.MarkReady()

	select {
	case <-metrics.CanStart:
		break
	case <-metrics.Done():
		metrics.MarkDone()
		return
	}

	log.Info().Msgf("Start metrics daemon, update each %v into %v", metrics.refreshRate, metrics.storage.Root)

	go func() {
		for {
			select {
			case <-metrics.Done():
				metrics.Persist()
				metrics.MarkDone()
				return
			case <-ticker.C:
				metrics.Persist()
			}
		}
	}()

	metrics.WaitStop()
	log.Info().Msg("Stop metrics daemon")
}
