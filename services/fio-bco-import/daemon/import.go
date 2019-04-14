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
	"strconv"
	"sync"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/config"
	"github.com/jancajthaml-openbank/fio-bco-import/http"
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"
	"github.com/jancajthaml-openbank/fio-bco-import/utils"

	localfs "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

// FioImport represents fio gateway to ledger-rest import subroutine
type FioImport struct {
	Support
	tenant        string
	fioGateway    string
	ledgerGateway string
	vaultGateway  string
	storage       *localfs.Storage
	refreshRate   time.Duration
	metrics       *Metrics
	system        *ActorSystem
	httpClient    http.Client
}

// NewFioImport returns fio import fascade
func NewFioImport(ctx context.Context, cfg config.Configuration, metrics *Metrics, system *ActorSystem, storage *localfs.Storage) FioImport {
	return FioImport{
		Support:       NewDaemonSupport(ctx),
		tenant:        cfg.Tenant,
		storage:       storage,
		fioGateway:    cfg.FioGateway,
		ledgerGateway: cfg.LedgerGateway,
		vaultGateway:  cfg.VaultGateway,
		refreshRate:   cfg.SyncRate,
		metrics:       metrics,
		system:        system,
		httpClient:    http.NewClient(),
	}
}

func (fio FioImport) getActiveTokens() ([]model.Token, error) {
	tokens, err := persistence.LoadTokens(fio.storage)
	if err != nil {
		return nil, err
	}
	uniq := make([]model.Token, 0, len(tokens))
	visited := make(map[string]bool)
	for _, token := range tokens {
		if _, ok := visited[token.Value]; !ok {
			visited[token.Value] = true
			uniq = append(uniq, token)
		}
	}
	return uniq, nil
}

func (fio FioImport) setLastSyncedID(token string, lastID int64) error {
	var (
		err      error
		response []byte
		code     int
		uri      string
	)

	if lastID != 0 {
		uri = fio.fioGateway + "/ib_api/rest/set-last-id/" + token + "/" + strconv.FormatInt(lastID, 10) + "/"
	} else {
		uri = fio.fioGateway + "/ib_api/rest/set-last-date/" + token + "/2012-07-27/"
	}

	response, code, err = fio.httpClient.Get(uri)
	if err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("FIO Gateway Error %d %+v", code, string(response))
	}

	return nil
}

func (fio FioImport) importNewTransactions(token model.Token) error {
	var (
		err      error
		request  []byte
		response []byte
		code     int
	)

	response, code, err = fio.httpClient.Get(fio.fioGateway + "/ib_api/rest/last/" + token.Value + "/transactions.json")
	if err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("fio gateway invalid response %d %+v", code, string(response))
	}

	var envelope model.FioImportEnvelope
	err = utils.JSON.Unmarshal(response, &envelope)

	if err != nil {
		return err
	}

	accounts := envelope.GetAccounts()

	for _, account := range accounts {
		request, err = utils.JSON.Marshal(account)
		if err != nil {
			return err
		}

		uri := fio.vaultGateway + "/account/" + fio.tenant
		err = utils.Retry(10, time.Second, func() (err error) {
			response, code, err = fio.httpClient.Post(uri, request)
			if code == 200 || code == 409 || code == 400 {
				return
			} else if code >= 500 && err == nil {
				err = fmt.Errorf("vault POST %s error %d %+v", uri, code, string(response))
			}
			return
		})

		if err != nil {
			return fmt.Errorf("vault account error %d %+v", code, err)
		} else if code == 400 {
			return fmt.Errorf("vault account malformed request %+v", string(request))
		} else if code != 200 && code != 409 {
			return fmt.Errorf("vault account error %d %+v", code, string(response))
		}
	}

	transactions := envelope.GetTransactions()

	var lastID int64

	for _, transaction := range transactions {

		for _, transfer := range transaction.Transfers {
			if transfer.IDTransfer > lastID {
				lastID = transfer.IDTransfer
			}
		}

		request, err = utils.JSON.Marshal(transaction)
		if err != nil {
			return err
		}

		uri := fio.ledgerGateway + "/transaction/" + fio.tenant
		err = utils.Retry(10, time.Second, func() (err error) {
			response, code, err = fio.httpClient.Post(uri, request)
			if code == 200 || code == 201 || code == 400 {
				return
			} else if code >= 500 && err == nil {
				err = fmt.Errorf("ledger-rest POST %s error %d %+v", uri, code, string(response))
			}
			return
		})

		if err != nil {
			return fmt.Errorf("ledger-rest transaction error %d %+v", code, err)
		} else if code == 409 {
			return fmt.Errorf("ledger-rest transaction duplicate %+v", string(request))
		} else if code == 400 {
			return fmt.Errorf("ledger-rest transaction malformed request %+v", string(request))
		} else if code != 200 && code != 201 {
			return fmt.Errorf("ledger-rest transaction error %d %+v", code, string(response))
		}

		if lastID != 0 {
			token.LastSyncedID = lastID
			if !persistence.UpdateToken(fio.storage, &token) {
				log.Warnf("Unable to update token %+v", token)
			}
		}

	}

	return nil
}

func (fio FioImport) importStatements(token model.Token) {
	if err := fio.setLastSyncedID(token.Value, token.LastSyncedID); err != nil {
		log.Warnf("set Last Synced ID Failed : %+v for %+v", err, token.ID)
		return
	}

	if fio.ctx.Err() != nil {
		return
	}

	if err := fio.importNewTransactions(token); err != nil {
		log.Warnf("import statements Failed : %+v for %+v", err, token.ID)
		return
	}
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

	var wg sync.WaitGroup

	for _, item := range tokens {
		wg.Add(1)
		go func(token model.Token) {
			defer wg.Done()
			fio.importStatements(token)
		}(item)
	}

	wg.Wait()
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
