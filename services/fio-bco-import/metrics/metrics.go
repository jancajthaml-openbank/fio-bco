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
	localfs "github.com/jancajthaml-openbank/local-fs"
	metrics "github.com/rcrowley/go-metrics"
)

// Metrics holds metrics counters
type Metrics struct {
	storage              localfs.Storage
	tenant               string
	continuous           bool
	createdTokens        metrics.Counter
	deletedTokens        metrics.Counter
	syncLatency          metrics.Timer
	importedTransfers    metrics.Meter
	importedTransactions metrics.Meter
}

// NewMetrics returns blank metrics holder
func NewMetrics(output string, continuous bool, tenant string) *Metrics {
	storage, err := localfs.NewPlaintextStorage(output)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	return &Metrics{
		continuous:           continuous,
		storage:              storage,
		tenant:               tenant,
		createdTokens:        metrics.NewCounter(),
		deletedTokens:        metrics.NewCounter(),
		syncLatency:          metrics.NewTimer(),
		importedTransfers:    metrics.NewMeter(),
		importedTransactions: metrics.NewMeter(),
	}
}

// TokenCreated increments token created by one
func (metrics *Metrics) TokenCreated() {
	if metrics == nil {
		return
	}
	metrics.createdTokens.Inc(1)
}

// TokenDeleted increments token deleted by one
func (metrics *Metrics) TokenDeleted() {
	if metrics == nil {
		return
	}
	metrics.deletedTokens.Inc(1)
}

// TimeSyncLatency measures time of sync duration
func (metrics *Metrics) TimeSyncLatency(f func()) {
	if metrics == nil {
		return
	}
	metrics.syncLatency.Time(f)
}

// TransactionImported increments transactions importer by one
func (metrics *Metrics) TransactionImported() {
	if metrics == nil {
		return
	}
	metrics.importedTransactions.Mark(1)
}

// TransfersImported increments transfers importer by given num
func (metrics *Metrics) TransfersImported(num int64) {
	if metrics == nil {
		return
	}
	metrics.importedTransfers.Mark(num)
}

// Setup hydrates metrics from storage
func (metrics *Metrics) Setup() error {
	if metrics == nil {
		return nil
	}
	if metrics.continuous {
		metrics.Hydrate()
	}
	return nil
}

// Done returns always finished
func (metrics *Metrics) Done() <-chan interface{} {
	done := make(chan interface{})
	close(done)
	return done
}

// Cancel does nothing
func (metrics *Metrics) Cancel() {
}

// Work represents metrics worker work
func (metrics *Metrics) Work() {
	metrics.Persist()
}
