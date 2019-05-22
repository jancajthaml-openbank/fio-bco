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

package actor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"
	"github.com/jancajthaml-openbank/fio-bco-import/utils"

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

// NilToken represents token that is neither existing neither non existing
func NilToken(s *ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		tokenHydration := persistence.LoadToken(s.Storage, state.ID)

		if tokenHydration == nil {
			context.Self.Become(state, NonExistToken(s))
			log.Debugf("%s ~ Nil -> NonExist", state.ID)
		} else {
			context.Self.Become(*tokenHydration, ExistToken(s))
			log.Debugf("%s ~ Nil -> Exist", state.ID)
		}

		context.Self.Receive(context)
	}
}

// NonExistToken represents token that does not exist
func NonExistToken(s *ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch msg := context.Data.(type) {

		case model.CreateToken:
			tokenResult := persistence.CreateToken(s.Storage, state.ID, msg.Value)

			if tokenResult == nil {
				s.SendRemote(FatalErrorMessage(context))
				log.Debugf("%s ~ (NonExist CreateToken) Error", state.ID)
				return
			}

			s.SendRemote(TokenCreatedMessage(context))
			log.Infof("New Token %s Created", state.ID)
			log.Debugf("%s ~ (NonExist CreateToken) OK", state.ID)
			s.Metrics.TokenCreated()

			context.Self.Become(*tokenResult, ExistToken(s))
			context.Self.Tell(model.SynchronizeToken{}, context.Receiver, context.Sender)

		case model.DeleteToken:
			s.SendRemote(FatalErrorMessage(context))
			log.Debugf("%s ~ (NonExist DeleteToken) Error", state.ID)

		case model.SynchronizeToken:
			break

		default:
			s.SendRemote(FatalErrorMessage(context))
			log.Debugf("%s ~ (NonExist Unknown Message) Error", state.ID)
		}

		return
	}
}

// ExistToken represents account that does exist
func ExistToken(s *ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch context.Data.(type) {

		case model.CreateToken:
			s.SendRemote(FatalErrorMessage(context))
			log.Debugf("%s ~ (Exist CreateToken) Error", state.ID)

		case model.SynchronizeToken:
			importStatements(s, state)
			log.Debugf("%s ~ (Exist SynchronizeToken) OK", state.ID)

		case model.DeleteToken:
			if !persistence.DeleteToken(s.Storage, state.ID) {
				s.SendRemote(FatalErrorMessage(context))
				log.Debugf("%s ~ (Exist DeleteToken) Error", state.ID)
				return
			}
			log.Infof("Token %s Deleted", state.ID)
			log.Debugf("%s ~ (Exist DeleteToken) OK", state.ID)
			s.Metrics.TokenDeleted()
			s.SendRemote(TokenDeletedMessage(context))
			context.Self.Become(state, NonExistToken(s))

		default:
			s.SendRemote(FatalErrorMessage(context))
			log.Warnf("%s ~ (Exist Unknown Message) Error", state.ID)

		}

		return
	}
}

func setLastSyncedID(s *ActorSystem, token model.Token) error {
	var (
		err      error
		response []byte
		code     int
		uri      string
	)

	if token.LastSyncedID != 0 {
		uri = s.FioGateway + "/ib_api/rest/set-last-id/" + token.Value + "/" + strconv.FormatInt(token.LastSyncedID, 10) + "/"
	} else {
		uri = s.FioGateway + "/ib_api/rest/set-last-date/" + token.Value + "/2012-07-27/"
	}

	response, code, err = s.HttpClient.Get(uri)
	if err != nil {
		return err
	}

	if code != 200 {
		return fmt.Errorf("fio gateway %s invalid response %d %+v", uri, code, string(response))
	}

	return nil
}

func importNewTransactions(s *ActorSystem, token model.Token) error {
	var (
		err      error
		request  []byte
		response []byte
		code     int
	)

	uri := s.FioGateway + "/ib_api/rest/last/" + token.Value + "/transactions.json"
	response, code, err = s.HttpClient.Get(uri)
	if err != nil {
		return err
	}

	if code != 200 {
		return fmt.Errorf("fio gateway %s invalid response %d %+v", uri, code, string(response))
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

		uri := s.VaultGateway + "/account/" + s.Tenant
		err = utils.Retry(10, time.Second, func() (err error) {
			response, code, err = s.HttpClient.Post(uri, request)
			if code == 200 || code == 409 || code == 400 {
				return
			} else if code >= 500 && err == nil {
				err = fmt.Errorf("vault-rest POST %s error %d %+v", uri, code, string(response))
			}
			return
		})

		if err != nil {
			return fmt.Errorf("vault-rest account error %d %+v", code, err)
		} else if code == 400 {
			return fmt.Errorf("vault-rest account malformed request %+v", string(request))
		} else if code != 200 && code != 409 {
			return fmt.Errorf("vault-rest account error %d %+v", code, string(response))
		}
	}

	transactions := envelope.GetTransactions(s.Tenant)

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

		uri := s.LedgerGateway + "/transaction/" + s.Tenant
		err = utils.Retry(10, time.Second, func() (err error) {
			response, code, err = s.HttpClient.Post(uri, request)
			if code == 200 || code == 201 || code == 400 {
				return
			} else if code >= 500 && err == nil {
				err = fmt.Errorf("ledger-rest POST %s error %d %+v", uri, code, string(response))
			}
			return
		})

		if err != nil {
			return fmt.Errorf("ledger-rest transaction error %d %+v", code, err)
		}
		if code == 409 {
			return fmt.Errorf("ledger-rest transaction duplicate %+v", string(request))
		}
		if code == 400 {
			return fmt.Errorf("ledger-rest transaction malformed request %+v", string(request))
		}
		if code != 200 && code != 201 {
			return fmt.Errorf("ledger-rest transaction error %d %+v", code, string(response))
		}

		if lastID != 0 {
			token.LastSyncedID = lastID
			if !persistence.UpdateToken(s.Storage, &token) {
				log.Warnf("Unable to update token %+v", token)
			}
		}

	}

	return nil
}

func importStatements(s *ActorSystem, token model.Token) {
	if err := setLastSyncedID(s, token); err != nil {
		log.Warnf("set Last Synced ID Failed : %+v for %+v", err, token.ID)
		return
	}

	if err := importNewTransactions(s, token); err != nil {
		log.Warnf("importNewTransactions failed %+v for %+v", err, token.ID)
		return
	}
}
