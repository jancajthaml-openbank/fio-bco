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

package integration

import (
	"context"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"
	"github.com/jancajthaml-openbank/fio-bco-import/utils"

	system "github.com/jancajthaml-openbank/actor-system"
	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// FioImport represents fio gateway to ledger-rest import subroutine
type FioImport struct {
	utils.DaemonSupport
	callback   func(msg interface{}, to system.Coordinates, from system.Coordinates)
	fioGateway string
	storage    *localfs.EncryptedStorage
	syncRate   time.Duration
}

// NewFioImport returns fio import fascade
func NewFioImport(ctx context.Context, fioEndpoint string, syncRate time.Duration, storage *localfs.EncryptedStorage, callback func(msg interface{}, to system.Coordinates, from system.Coordinates)) FioImport {
	return FioImport{
		DaemonSupport: utils.NewDaemonSupport(ctx, "fio"),
		callback:      callback,
		storage:       storage,
		fioGateway:    fioEndpoint,
		syncRate:      syncRate,
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

	for _, item := range tokens {
		log.Debugf("Request to import token %s", item)
		msg := model.SynchronizeToken{}
		to := system.Coordinates{Name: item}
		from := system.Coordinates{Name: "token_import_cron"}
		fio.callback(msg, to, from)
	}
}

// Start handles everything needed to start fio import daemon
func (fio FioImport) Start() {
	fio.MarkReady()

	select {
	case <-fio.CanStart:
		break
	case <-fio.Done():
		fio.MarkDone()
		return
	}

	log.Infof("Start fio-import daemon, sync %v now and then each %v", fio.fioGateway, fio.syncRate)

	fio.importRoundtrip()

	go func() {
		for {
			select {
			case <-fio.Done():
				fio.MarkDone()
				return
			case <-time.After(fio.syncRate):
				fio.importRoundtrip()
			}
		}
	}()

	fio.WaitStop()
	log.Info("Stop fio-import daemon")
}
