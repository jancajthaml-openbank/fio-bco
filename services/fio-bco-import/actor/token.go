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

	"github.com/jancajthaml-openbank/fio-bco-import/fio"
	"github.com/jancajthaml-openbank/fio-bco-import/ledger"
	"github.com/jancajthaml-openbank/fio-bco-import/metrics"
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"
	"github.com/jancajthaml-openbank/fio-bco-import/vault"

	system "github.com/jancajthaml-openbank/actor-system"
)

// NilToken represents token that is neither existing neither non existing
func NilToken(s *ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		tokenHydration := persistence.LoadToken(s.Storage, state.ID)

		if tokenHydration == nil {
			context.Self.Become(state, NonExistToken(s))
			log.WithField("token", state.ID).Debug("Nil -> NonExist")
		} else {
			context.Self.Become(*tokenHydration, ExistToken(s))
			log.WithField("token", state.ID).Debug("Nil -> Exist")
		}

		context.Self.Receive(context)
	}
}

// NonExistToken represents token that does not exist
func NonExistToken(s *ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch msg := context.Data.(type) {

		case ProbeMessage:
			break

		case CreateToken:
			tokenResult := persistence.CreateToken(s.Storage, state.ID, msg.Value)

			if tokenResult == nil {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.WithField("token", state.ID).Debug("(NonExist CreateToken) Error")
				return
			}

			s.SendMessage(RespCreateToken, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Info("New Token Created")
			log.WithField("token", state.ID).Debug("(NonExist CreateToken) OK")
			s.Metrics.TokenCreated()

			context.Self.Become(*tokenResult, ExistToken(s))
			context.Self.Tell(SynchronizeToken{}, context.Receiver, context.Sender)

		case DeleteToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Debug("(NonExist DeleteToken) Error")

		case SynchronizeToken:
			break

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Debug("(NonExist Unknown Message) Error")
		}

		return
	}
}

// ExistToken represents account that does exist
func ExistToken(s *ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch context.Data.(type) {

		case ProbeMessage:
			break

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Debug("(Exist CreateToken) Error")

		case SynchronizeToken:
			log.WithField("token", state.ID).Debug("(Exist SynchronizeToken)")
			context.Self.Become(t_state, SynchronizingToken(s))
			go importStatements(s, state, func() {
				context.Self.Become(t_state, NilToken(s))
				context.Self.Tell(ProbeMessage{}, context.Receiver, context.Receiver)
			})

		case DeleteToken:
			if !persistence.DeleteToken(s.Storage, state.ID) {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.WithField("token", state.ID).Debugf("(Exist DeleteToken) Error")
				return
			}
			log.WithField("token", state.ID).Info("Token Deleted")
			log.WithField("token", state.ID).Debug("(Exist DeleteToken) OK")
			s.Metrics.TokenDeleted()
			s.SendMessage(RespDeleteToken, context.Sender, context.Receiver)
			context.Self.Become(state, NonExistToken(s))

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Warn("(Exist Unknown Message) Error")

		}

		return
	}
}

// SynchronizingToken represents account that is currently synchronizing
func SynchronizingToken(s *ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {
		state := t_state.(model.Token)

		switch context.Data.(type) {

		case ProbeMessage:
			break

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Debug("(Synchronizing CreateToken) Error")

		case SynchronizeToken:
			log.WithField("token", state.ID).Debug("(Synchronizing SynchronizeToken)")

		case DeleteToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Debug("(Synchronizing DeleteToken) Error")

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.WithField("token", state.ID).Warn("(Synchronizing Unknown Message) Error")

		}

		return
	}
}

func importNewStatements(tenant string, fioClient *fio.FioClient, vaultClient *vault.VaultClient, ledgerClient *ledger.LedgerClient, metrics *metrics.Metrics, token *model.Token) (int64, error) {
	var (
		statements *fio.FioImportEnvelope
		err        error
		lastID     int64 = token.LastSyncedID
	)

	metrics.TimeSyncLatency(func() {
		statements, err = fioClient.GetTransactions()
	})
	if err != nil {
		return lastID, err
	}
	if len(statements.Statement.TransactionList.Transactions) == 0 {
		return lastID, nil
	}

	// FIXME getStatements end here

	log.WithField("token", token.ID).Debugf("sorting statements")

	sort.SliceStable(statements.Statement.TransactionList.Transactions, func(i, j int) bool {
		return statements.Statement.TransactionList.Transactions[i].TransactionID.Value == statements.Statement.TransactionList.Transactions[j].TransactionID.Value
	})

	log.WithField("token", token.ID).Debugf("importing accounts")

	for account := range statements.GetAccounts(tenant) {
		err = vaultClient.CreateAccount(account)
		if err != nil {
			return lastID, err
		}
	}

	log.WithField("token", token.ID).Debugf("importing transactions")

	for transaction := range statements.GetTransactions(tenant) {
		err = ledgerClient.CreateTransaction(transaction)
		if err != nil {
			return lastID, err
		}

		metrics.TransactionImported()
		metrics.TransfersImported(int64(len(transaction.Transfers)))

		for _, transfer := range transaction.Transfers {
			if transfer.IDTransfer > lastID {
				lastID = transfer.IDTransfer
			}
		}
	}

	return lastID, nil
}

func importStatements(s *ActorSystem, token model.Token, callback func()) {
	defer callback()

	log.WithField("token", token.ID).Debug("Importing statements")

	fioClient := fio.NewFioClient(s.FioGateway, token)
	vaultClient := vault.NewVaultClient(s.VaultGateway)
	ledgerClient := ledger.NewLedgerClient(s.LedgerGateway)

	log.WithField("token", token.ID).Debug("Import Begin")
	lastID, err := importNewStatements(s.Tenant, &fioClient, &vaultClient, &ledgerClient, s.Metrics, &token)
	if err != nil {
		log.WithField("token", token.ID).Errorf("Import statements failed with %+v", err)
	}
	if lastID > token.LastSyncedID {
		token.LastSyncedID = lastID
		if !persistence.UpdateToken(s.Storage, &token) {
			log.WithField("token", token.ID).Warn("Unable to update token")
		}
	}
	log.WithField("token", token.ID).Debug("Import End")
}
