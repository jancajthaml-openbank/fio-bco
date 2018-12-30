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

package daemon

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-unit/config"
	"github.com/jancajthaml-openbank/fio-bco-unit/http"
	"github.com/jancajthaml-openbank/fio-bco-unit/model"
	"github.com/jancajthaml-openbank/fio-bco-unit/persistence"
	"github.com/jancajthaml-openbank/fio-bco-unit/utils"

	log "github.com/sirupsen/logrus"
)

type FioImport struct {
	Support
	tenant      string
	fioGateway  string
	wallGateway string
	storage     string
	refreshRate time.Duration
	metrics     *Metrics
	system      *ActorSystem
	httpClient  http.Client
}

func NewFioImport(ctx context.Context, cfg config.Configuration, metrics *Metrics, system *ActorSystem) FioImport {
	return FioImport{
		Support:     NewDaemonSupport(ctx),
		tenant:      cfg.Tenant,
		storage:     cfg.RootStorage,
		fioGateway:  cfg.FioGateway,
		wallGateway: cfg.WallGateway,
		refreshRate: cfg.SyncRate,
		metrics:     metrics,
		system:      system,
		httpClient:  http.NewClient(),
	}
}

func (fio FioImport) getActiveTokens() ([]model.Token, error) {
	return persistence.LoadTokens(fio.storage)
}

func (fio FioImport) setLastSyncedID(token string, lastID int64) error {
	var (
		err  error
		code int
		uri  string
	)

	if lastID != 0 {
		uri = fio.fioGateway + "/ib_api/rest/set-last-id/" + token + "/" + strconv.FormatInt(lastID, 10)
	} else {
		uri = fio.fioGateway + "/ib_api/rest/set-last-date/" + token + "/2012-07-27"
	}

	_, code, err = fio.httpClient.Get(uri)
	if err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("%d", code)
	}

	return nil
}

func (fio FioImport) importNewTransactions(token model.Token) error {
	var (
		err  error
		data []byte
		code int
	)

	data, code, err = fio.httpClient.Get(fio.fioGateway + "/ib_api/rest/last/" + token.Value + "/transactions.json")
	if err != nil {
		return err
	} else if code != 200 {
		return fmt.Errorf("%d", code)
	}

	var envelope model.FioImportEnvelope
	err = utils.JSON.Unmarshal(data, &envelope)

	if err != nil {
		return err
	}

	accounts := envelope.GetAccounts()

	for _, account := range accounts {
		data, err = utils.JSON.Marshal(account)
		if err != nil {
			return err
		}

		err = utils.Retry(10, time.Second, func() (err error) {
			data, code, err = fio.httpClient.Post(fio.wallGateway+"/account/"+fio.tenant, data)
			if code == 200 || code == 409 {
				return
			} else if code >= 500 && err == nil {
				err = fmt.Errorf("%d", code)
			}
			return
		})

		if err != nil {
			return err
		} else if code != 200 && code != 409 {
			return fmt.Errorf("%d", code)
		}
	}

	transactions := envelope.GetTransactions()

	var lastId int64 = 0

	for _, transaction := range transactions {

		for _, transfer := range transaction.Transfers {
			if transfer.IDTransfer > lastId {
				lastId = transfer.IDTransfer
			}
		}

		data, err = utils.JSON.Marshal(transaction)
		if err != nil {
			return err
		}

		err = utils.Retry(10, time.Second, func() (err error) {
			data, code, err = fio.httpClient.Post(fio.wallGateway+"/transaction/"+fio.tenant, data)
			if code == 200 || code == 201 {
				return
			} else if code >= 500 && err == nil {
				err = fmt.Errorf("%d", code)
			}
			return
		})

		if err != nil {
			return err
		} else if code != 200 && code != 201 {
			return fmt.Errorf("%d", code)
		}

		if lastId != 0 {
			token.LastSyncedID = lastId
			if !persistence.UpdateToken(fio.storage, &token) {
				log.Warnf("Unable to update token %+v", token)
			}
		}

	}

	return nil
}

func (fio FioImport) importStatements(token model.Token) {
	if err := fio.setLastSyncedID(token.Value, token.LastSyncedID); err != nil {
		log.Warnf("FIO Gateway returned error %+v for %+v", err, token.Value)
		return
	}

	if err := fio.importNewTransactions(token); err != nil {
		log.Warnf("FIO Gateway returned error %+v for %+v", err, token.Value)
		return
	}
}

func (fio FioImport) importRoundtrip() {
	var wg sync.WaitGroup

	tokens, err := fio.getActiveTokens()
	if err != nil {
		log.Errorf("Unable to get active tokens %+v", err)
		return
	}

	for _, item := range tokens {
		wg.Add(1)
		go func(token model.Token) {
			defer wg.Done()
			fio.importStatements(token)
		}(item)
	}

	wg.Wait()
}

func (fio FioImport) Start() {
	defer fio.MarkDone()

	log.Infof("Start fio-import daemon, sync %v each %v", fio.fioGateway, fio.refreshRate)

	fio.MarkReady()

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
