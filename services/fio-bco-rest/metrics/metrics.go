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
	"time"

	"github.com/jancajthaml-openbank/fio-bco-rest/utils"

	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

// Metrics represents metrics subroutine
type Metrics struct {
	utils.DaemonSupport
	output             string
	refreshRate        time.Duration
	getTokenLatency    metrics.Timer
	createTokenLatency metrics.Timer
	deleteTokenLatency metrics.Timer
}

// NewMetrics returns metrics fascade
func NewMetrics(ctx context.Context, output string, refreshRate time.Duration) Metrics {
	return Metrics{
		DaemonSupport:      utils.NewDaemonSupport(ctx),
		output:             output,
		refreshRate:        refreshRate,
		createTokenLatency: metrics.NewTimer(),
		deleteTokenLatency: metrics.NewTimer(),
		getTokenLatency:    metrics.NewTimer(),
	}
}

// TimeGetToken measure execution of GetToken
func (metrics *Metrics) TimeGetToken(f func()) {
	metrics.getTokenLatency.Time(f)
}

// TimeCreateToken measure execution of CreateToken
func (metrics *Metrics) TimeCreateToken(f func()) {
	metrics.createTokenLatency.Time(f)
}

// TimeDeleteToken measure execution of DeleteToken
func (metrics *Metrics) TimeDeleteToken(f func()) {
	metrics.deleteTokenLatency.Time(f)
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
