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
	storage            localfs.Storage
	continuous         bool
	getTokenLatency    metrics.Timer
	createTokenLatency metrics.Timer
	deleteTokenLatency metrics.Timer
}

// NewMetrics returns blank metrics holder
func NewMetrics(output string, continuous bool) *Metrics {
	storage, err := localfs.NewPlaintextStorage(output)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	return &Metrics{
		storage:            storage,
		continuous:         continuous,
		createTokenLatency: metrics.NewTimer(),
		deleteTokenLatency: metrics.NewTimer(),
		getTokenLatency:    metrics.NewTimer(),
	}
}

// TimeGetToken measure execution of GetToken
func (metrics *Metrics) TimeGetToken(f func()) {
	if metrics == nil {
		return
	}
	metrics.getTokenLatency.Time(f)
}

// TimeCreateToken measure execution of CreateToken
func (metrics *Metrics) TimeCreateToken(f func()) {
	if metrics == nil {
		return
	}
	metrics.createTokenLatency.Time(f)
}

// TimeDeleteToken measure execution of DeleteToken
func (metrics *Metrics) TimeDeleteToken(f func()) {
	if metrics == nil {
		return
	}
	metrics.deleteTokenLatency.Time(f)
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
