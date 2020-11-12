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

	"github.com/jancajthaml-openbank/fio-bco-rest/utils"
	localfs "github.com/jancajthaml-openbank/local-fs"
	metrics "github.com/rcrowley/go-metrics"
)

// Metrics holds metrics counters
type Metrics struct {
	utils.DaemonSupport
	storage            localfs.Storage
	refreshRate        time.Duration
	getTokenLatency    metrics.Timer
	createTokenLatency metrics.Timer
	deleteTokenLatency metrics.Timer
}

// NewMetrics returns blank metrics holder
func NewMetrics(ctx context.Context, output string, refreshRate time.Duration) *Metrics {
	storage, err := localfs.NewPlaintextStorage(output)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	return &Metrics{
		DaemonSupport:      utils.NewDaemonSupport(ctx, "metrics"),
		storage:            storage,
		refreshRate:        refreshRate,
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

// Start handles everything needed to start metrics daemon
func (metrics *Metrics) Start() {
	if metrics == nil {
		return
	}

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

	log.Info().Msgf("Start metrics daemon, update file each %v", metrics.refreshRate)

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
