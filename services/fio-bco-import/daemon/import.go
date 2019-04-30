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

package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/config"
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"

	system "github.com/jancajthaml-openbank/actor-system"
	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// FioImport represents fio gateway to ledger-rest import subroutine
type FioImport struct {
	Support
	callback    func(msg interface{}, to system.Coordinates, from system.Coordinates)
	fioGateway  string
	storage     *localfs.Storage
	refreshRate time.Duration
}

// NewFioImport returns fio import fascade
func NewFioImport(ctx context.Context, cfg config.Configuration, storage *localfs.Storage, callback func(msg interface{}, to system.Coordinates, from system.Coordinates)) FioImport {
	return FioImport{
		Support:     NewDaemonSupport(ctx),
		callback:    callback,
		storage:     storage,
		fioGateway:  cfg.FioGateway,
		refreshRate: cfg.SyncRate,
	}
}

func (fio FioImport) getActiveTokens() ([]string, error) {
	tokens, err := persistence.LoadTokens(fio.storage)
	if err != nil {
		return nil, err
	}
	notBefore := time.Now().Add(time.Second * time.Duration(-6))
	uniq := make([]string, 0)
	visited := make(map[string]bool)
	for _, token := range tokens {
		if !token.CreatedAt.Before(notBefore) {
			continue
		}
		if _, ok := visited[token.Value]; !ok {
			visited[token.Value] = true
			uniq = append(uniq, token.ID)
		}
	}
	return uniq, nil
}

func (fio FioImport) importRoundtrip() {
	tokens, err := fio.getActiveTokens()
	if err != nil {
		log.Errorf("unable to get active tokens %+v", err)
		return
	}

	if fio.ctx.Err() != nil {
		return
	}

	for _, item := range tokens {
		log.Debugf("Request to import token %s", item)
		msg := model.SynchronizeToken{}
		to := system.Coordinates{Name: item}
		from := system.Coordinates{Name: "token_import_cron"}
		fio.callback(msg, to, from)
	}
}

// WaitReady wait for fio import to be ready
func (fio FioImport) WaitReady(deadline time.Duration) (err error) {
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
	case <-fio.IsReady:
		ticker.Stop()
		err = nil
		return
	case <-ticker.C:
		err = fmt.Errorf("daemon was not ready within %v seconds", deadline)
		return
	}
}

// Start handles everything needed to start fio import daemon
func (fio FioImport) Start() {
	defer fio.MarkDone()

	fio.MarkReady()

	select {
	case <-fio.canStart:
		break
	case <-fio.Done():
		return
	}

	log.Infof("Start fio-import daemon, sync %v now and then each %v", fio.fioGateway, fio.refreshRate)

	fio.importRoundtrip()

	for {
		select {
		case <-fio.Done():
			log.Info("Stopping fio-import daemon")
			log.Info("Stop fio-import daemon")
			return
		case <-time.After(fio.refreshRate):
			fio.importRoundtrip()
		}
	}
}
