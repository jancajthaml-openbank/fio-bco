// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"sync/atomic"

	"github.com/DataDog/datadog-go/statsd"
)

// Metrics provides helper function for metrics
type Metrics interface {
	TokenCreated()
	TokenDeleted()
	TransactionImported(transfers int)
}

// StatsdMetrics provides metrics helper with statsd client
type StatsdMetrics struct {
	client               *statsd.Client
	tenant               string
	createdTokens        int64
	deletedTokens        int64
	importedTransactions int64
	importedTransfers    int64
}

// NewMetrics returns blank metrics holder
func NewMetrics(tenant string, endpoint string) *StatsdMetrics {
	client, err := statsd.New(endpoint, statsd.WithClientSideAggregation(), statsd.WithoutTelemetry())
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure statsd client")
		return nil
	}
	return &StatsdMetrics{
		client:               client,
		tenant:               tenant,
		createdTokens:        int64(0),
		deletedTokens:        int64(0),
		importedTransactions: int64(0),
		importedTransfers:    int64(0),
	}
}

// TokenCreated increments token created by one
func (instance *StatsdMetrics) TokenCreated() {
	if instance == nil {
		return
	}
	atomic.AddInt64(&(instance.createdTokens), 1)
}

// TokenDeleted increments token deleted by one
func (instance *StatsdMetrics) TokenDeleted() {
	if instance == nil {
		return
	}
	//metrics.deletedTokens.Inc(1)
	atomic.AddInt64(&(instance.deletedTokens), 1)
}

// TransactionImported increments transactions importer by one
func (instance *StatsdMetrics) TransactionImported(transfers int) {
	if instance == nil {
		return
	}
	atomic.AddInt64(&(instance.importedTransactions), 1)
	atomic.AddInt64(&(instance.importedTransfers), int64(transfers))
}

// Setup does nothing
func (*StatsdMetrics) Setup() error {
	return nil
}

// Done returns always finished
func (*StatsdMetrics) Done() <-chan interface{} {
	done := make(chan interface{})
	close(done)
	return done
}

// Cancel does nothing
func (*StatsdMetrics) Cancel() {
}

// Work represents metrics worker work
func (instance *StatsdMetrics) Work() {
	if instance == nil {
		return
	}

	createdTokens := instance.createdTokens
	deletedTokens := instance.deletedTokens
	importedTransactions := instance.importedTransactions
	importedTransfers := instance.importedTransfers

	atomic.AddInt64(&(instance.createdTokens), -createdTokens)
	atomic.AddInt64(&(instance.deletedTokens), -deletedTokens)
	atomic.AddInt64(&(instance.importedTransactions), -importedTransactions)
	atomic.AddInt64(&(instance.importedTransfers), -importedTransfers)

	tags := []string{"tenant:" + instance.tenant}

	instance.client.Count("openbank.bco.fio.token.created", createdTokens, tags, 1)
	instance.client.Count("openbank.bco.fio.token.deleted", deletedTokens, tags, 1)
	instance.client.Count("openbank.bco.fio.transaction.imported", importedTransactions, tags, 1)
	instance.client.Count("openbank.bco.fio.transfer.imported", importedTransfers, tags, 1)
}
