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

package boot

import (
	"context"
	"os"

	"github.com/jancajthaml-openbank/fio-bco-rest/actor"
	"github.com/jancajthaml-openbank/fio-bco-rest/api"
	"github.com/jancajthaml-openbank/fio-bco-rest/config"
	"github.com/jancajthaml-openbank/fio-bco-rest/daemon"
	"github.com/jancajthaml-openbank/fio-bco-rest/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// Application encapsulate initialized application
type Application struct {
	cfg           config.Configuration
	interrupt     chan os.Signal
	actorSystem   daemon.ActorSystem
	rest          daemon.Server
	systemControl daemon.SystemControl
	cancel        context.CancelFunc
}

// Initialize application
func Initialize() Application {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.GetConfig()

	utils.SetupLogger(cfg.LogLevel)

	log.Infof(">>> Setup <<<")

	systemControl := daemon.NewSystemControl(ctx, cfg)

	storage := localfs.NewStorage(cfg.RootStorage)
	storage.SetEncryptionKey(cfg.EncryptionKey)

	actorSystem := daemon.NewActorSystem(ctx, cfg)
	actorSystem.Support.RegisterOnRemoteMessage(actor.ProcessRemoteMessage(&actorSystem))

	rest := daemon.NewServer(ctx, cfg)
	rest.HandleFunc("/health", api.HealtCheck, "GET", "HEAD")
	rest.HandleFunc("/tenant/{tenant}", api.TenantPartial(&systemControl), "POST", "DELETE")
	rest.HandleFunc("/tenant", api.TenantsPartial(&systemControl), "GET")
	rest.HandleFunc("/token/{tenant}/{token}", api.TokenPartial(&actorSystem), "DELETE")
	rest.HandleFunc("/token/{tenant}", api.TokensPartial(&actorSystem, &storage), "POST", "GET")

	return Application{
		cfg:           cfg,
		interrupt:     make(chan os.Signal, 1),
		actorSystem:   actorSystem,
		rest:          rest,
		systemControl: systemControl,
		cancel:        cancel,
	}
}
