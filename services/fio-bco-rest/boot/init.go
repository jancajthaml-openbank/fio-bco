// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/jancajthaml-openbank/fio-bco-rest/actor"
	"github.com/jancajthaml-openbank/fio-bco-rest/api"
	"github.com/jancajthaml-openbank/fio-bco-rest/config"
	"github.com/jancajthaml-openbank/fio-bco-rest/daemon"
	"github.com/jancajthaml-openbank/fio-bco-rest/utils"

	log "github.com/sirupsen/logrus"
)

// Application encapsulate initialized application
type Application struct {
	cfg         config.Configuration
	interrupt   chan os.Signal
	actorSystem daemon.ActorSystem
	rest        daemon.Server
	cancel      context.CancelFunc
}

// Initialize application
func Initialize() Application {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewResolver().GetConfig()

	log.SetFormatter(new(utils.LogFormat))

	log.Infof(">>> Setup <<<")

	if cfg.LogOutput == "" {
		log.SetOutput(os.Stdout)
	} else if file, err := os.OpenFile(cfg.LogOutput, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600); err == nil {
		defer file.Close()
		log.SetOutput(bufio.NewWriter(file))
	} else {
		log.SetOutput(os.Stdout)
		log.Warnf("Unable to create %s: %v", cfg.LogOutput, err)
	}

	if level, err := log.ParseLevel(cfg.LogLevel); err == nil {
		log.Infof("Log level set to %v", strings.ToUpper(cfg.LogLevel))
		log.SetLevel(level)
	} else {
		log.Warnf("Invalid log level %v, using level WARN", cfg.LogLevel)
		log.SetLevel(log.WarnLevel)
	}

	actorSystem := daemon.NewActorSystem(ctx, cfg)
	actorSystem.Support.RegisterOnRemoteMessage(actor.ProcessRemoteMessage(&actorSystem))

	rest := daemon.NewServer(ctx, cfg, &actorSystem)
	rest.HandleFunc("/health", api.HealtCheck, "GET", "HEAD")
	rest.HandleFunc("/token/{tenant_id}", api.CreateTokenPartial(&actorSystem), "POST")
	rest.HandleFunc("/token/{tenant_id}/{token_id}", api.DeleteTokenPartial(&actorSystem), "DELETE")
	rest.HandleFunc("/tokens/{tenant_id}", api.GetTokensPartial(cfg, &actorSystem), "GET")

	return Application{
		cfg:         cfg,
		interrupt:   make(chan os.Signal, 1),
		actorSystem: actorSystem,
		rest:        rest,
		cancel:      cancel,
	}
}
