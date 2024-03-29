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

package boot

import (
	"os"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/actor"
	"github.com/jancajthaml-openbank/fio-bco-import/config"
	"github.com/jancajthaml-openbank/fio-bco-import/integration"
	"github.com/jancajthaml-openbank/fio-bco-import/metrics"
	"github.com/jancajthaml-openbank/fio-bco-import/support/concurrent"
	"github.com/jancajthaml-openbank/fio-bco-import/support/logging"

	system "github.com/jancajthaml-openbank/actor-system"
)

// Program encapsulate program
type Program struct {
	interrupt chan os.Signal
	cfg       config.Configuration
	pool      concurrent.DaemonPool
}

// NewProgram returns new program
func NewProgram() Program {
	return Program{
		interrupt: make(chan os.Signal, 1),
		cfg:       config.LoadConfig(),
		pool:      concurrent.NewDaemonPool("program"),
	}
}

// Setup setups program
func (prog *Program) Setup() {
	if prog == nil {
		return
	}

	logging.SetupLogger(prog.cfg.LogLevel)

	metricsWorker := metrics.NewMetrics(
		prog.cfg.Tenant,
		prog.cfg.MetricsStastdEndpoint,
	)

	actorSystem := actor.NewActorSystem(
		prog.cfg.Tenant,
		prog.cfg.LakeHostname,
		prog.cfg.FioGateway,
		prog.cfg.VaultGateway,
		prog.cfg.LedgerGateway,
		prog.cfg.RootStorage,
		prog.cfg.EncryptionKey,
		metricsWorker,
	)

	fioWorker := integration.NewFioImport(
		prog.cfg.RootStorage,
		prog.cfg.EncryptionKey,
		func(token string) {
			actorSystem.SendMessage(actor.ReqSynchronizeToken,
				system.Coordinates{
					Region: actorSystem.Name,
					Name:   token,
				},
				system.Coordinates{
					Region: actorSystem.Name,
					Name:   "",
				},
			)
		},
	)

	prog.pool.Register(concurrent.NewOneShotDaemon(
		"actor-system",
		actorSystem,
	))

	prog.pool.Register(concurrent.NewScheduledDaemon(
		"metrics",
		metricsWorker,
		time.Second,
	))

	prog.pool.Register(concurrent.NewScheduledDaemon(
		"fio",
		fioWorker,
		prog.cfg.SyncRate,
	))
}
