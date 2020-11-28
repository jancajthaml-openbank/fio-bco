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

package integration

import (
	"context"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/persistence"
	"github.com/jancajthaml-openbank/fio-bco-import/support/concurrent"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// FioImport represents fio gateway to ledger-rest import subroutine
type FioImport struct {
	concurrent.DaemonSupport
	callback   func(token string)
	fioGateway string
	storage    localfs.Storage
	syncRate   time.Duration
}

// NewFioImport returns fio import fascade
func NewFioImport(ctx context.Context, fioEndpoint string, syncRate time.Duration, rootStorage string, storageKey []byte, callback func(token string)) *FioImport {
	storage, err := localfs.NewEncryptedStorage(rootStorage, storageKey)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	return &FioImport{
		DaemonSupport: concurrent.NewDaemonSupport(ctx, "fio"),
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
		log.Error().Msgf("unable to get active tokens %+v", err)
		return
	}

	for _, item := range tokens {
		log.Debug().Msgf("Request to import token %s", item)
		fio.callback(item)
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

	log.Info().Msgf("Start fio-import daemon, sync %v now and then each %v", fio.fioGateway, fio.syncRate)

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
	log.Info().Msg("Stop fio-import daemon")
}
