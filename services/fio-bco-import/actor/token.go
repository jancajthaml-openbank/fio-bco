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

package actor

import (
	"sort"

	"github.com/jancajthaml-openbank/fio-bco-import/metrics"
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"
	"github.com/jancajthaml-openbank/fio-bco-import/support/http"

	system "github.com/jancajthaml-openbank/actor-system"
)

// NilToken represents token that is neither existing neither non existing
func NilToken(s *System) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		tokenHydration := persistence.LoadToken(s.Storage, state.ID)

		if tokenHydration == nil {
			context.Self.Become(state, NonExistToken(s))
			log.Debug().Msgf("token %s Nil -> NonExist", state.ID)
		} else {
			context.Self.Become(*tokenHydration, ExistToken(s))
			log.Debug().Msgf("token %s Nil -> Exist", state.ID)
		}

		context.Self.Receive(context)
	}
}

// NonExistToken represents token that does not exist
func NonExistToken(s *System) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch msg := context.Data.(type) {

		case ProbeMessage:
			break

		case CreateToken:
			tokenResult := persistence.CreateToken(s.Storage, state.ID, msg.Value)

			if tokenResult == nil {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (NonExist CreateToken) Error", state.ID)
				return
			}

			s.SendMessage(RespCreateToken, context.Sender, context.Receiver)
			log.Info().Msgf("New Token %s Created", state.ID)
			log.Debug().Msgf("token %s (NonExist CreateToken) OK", state.ID)
			s.Metrics.TokenCreated()

			context.Self.Become(*tokenResult, ExistToken(s))
			context.Self.Tell(SynchronizeToken{}, context.Receiver, context.Sender)

		case DeleteToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist DeleteToken) Error", state.ID)

		case SynchronizeToken:
			break

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist Unknown Message) Error", state.ID)
		}

		return
	}
}

// ExistToken represents account that does exist
func ExistToken(s *System) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch context.Data.(type) {

		case ProbeMessage:
			break

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Exist CreateToken) Error", state.ID)

		case SynchronizeToken:
			log.Debug().Msgf("token %s (Exist SynchronizeToken)", state.ID)
			context.Self.Become(t_state, SynchronizingToken(s))
			go importStatements(s, state, func() {
				context.Self.Become(t_state, NilToken(s))
				context.Self.Tell(ProbeMessage{}, context.Receiver, context.Receiver)
			})

		case DeleteToken:
			if !persistence.DeleteToken(s.Storage, state.ID) {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (Exist DeleteToken) Error", state.ID)
				return
			}

			log.Info().Msgf("Token %s Deleted", state.ID)
			log.Debug().Msgf("token %s (Exist DeleteToken) OK", state.ID)
			s.Metrics.TokenDeleted()
			s.SendMessage(RespDeleteToken, context.Sender, context.Receiver)
			context.Self.Become(state, NonExistToken(s))

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Exist Unknown Message) Error", state.ID)

		}

		return
	}
}

// SynchronizingToken represents account that is currently synchronizing
func SynchronizingToken(s *System) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch context.Data.(type) {

		case ProbeMessage:
			break

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Synchronizing CreateToken) Error", state.ID)

		case SynchronizeToken:
			log.Debug().Msgf("token %s (Synchronizing SynchronizeToken)", state.ID)

		case DeleteToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Synchronizing DeleteToken) Error", state.ID)

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Warn().Msgf("token %s (Synchronizing Unknown Message) Error", state.ID)

		}

		return
	}
}

func importNewStatements(tenant string, fioClient *http.FioClient, vaultClient *http.VaultClient, ledgerClient *http.LedgerClient, metrics metrics.Metrics, token *model.Token) (int64, error) {
	var (
		statements *model.ImportEnvelope
		err        error
		lastID     int64 = token.LastSyncedID
	)

	statements, err = fioClient.GetTransactions()
	if err != nil {
		return lastID, err
	}
	if len(statements.Statement.TransactionList.Transactions) == 0 {
		return lastID, nil
	}

	// FIXME getStatements end here

	log.Debug().Msgf("token %s sorting statements", token.ID)

	sort.SliceStable(statements.Statement.TransactionList.Transactions, func(i, j int) bool {
		return statements.Statement.TransactionList.Transactions[i].TransactionID.Value == statements.Statement.TransactionList.Transactions[j].TransactionID.Value
	})

	log.Debug().Msgf("token %s importing accounts", token.ID)

	for account := range statements.GetAccounts(tenant) {
		err = vaultClient.CreateAccount(account)
		if err != nil {
			return lastID, err
		}
	}

	log.Debug().Msgf("token %s importing transactions", token.ID)

	for transaction := range statements.GetTransactions(tenant) {
		err = ledgerClient.CreateTransaction(transaction)
		if err != nil {
			return lastID, err
		}

		metrics.TransactionImported(len(transaction.Transfers))

		for _, transfer := range transaction.Transfers {
			if transfer.IDTransfer > lastID {
				lastID = transfer.IDTransfer
			}
		}
	}

	return lastID, nil
}

func importStatements(s *System, token model.Token, callback func()) {
	defer callback()

	log.Debug().Msgf("token %s Importing statements", token.ID)

	fioClient := http.NewFioClient(s.FioGateway, token)
	vaultClient := http.NewVaultClient(s.VaultGateway)
	ledgerClient := http.NewLedgerClient(s.LedgerGateway)

	log.Debug().Msgf("token %s Import Begin", token.ID)
	lastID, err := importNewStatements(s.Tenant, &fioClient, &vaultClient, &ledgerClient, s.Metrics, &token)
	if err != nil {
		log.Error().Msgf("token %s Import statements failed with %+v", token.ID, err)
	}
	if lastID > token.LastSyncedID {
		token.LastSyncedID = lastID
		if !persistence.UpdateToken(s.Storage, &token) {
			log.Warn().Msgf("token %s Unable to update token", token.ID)
		}
	}
	log.Debug().Msgf("token %s Import End", token.ID)
}
