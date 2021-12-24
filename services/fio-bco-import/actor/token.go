// Copyright (c) 2016-2021, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"github.com/jancajthaml-openbank/fio-bco-import/integration"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"

	system "github.com/jancajthaml-openbank/actor-system"
)

// NilToken represents token that is neither existing neither non existing
func NilToken(s *System, id string) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {
		context.Self.Tell(context.Data, context.Receiver, context.Sender)
		_, err := persistence.LoadToken(s.EncryptedStorage, id)
		if err != nil {
			log.Debug().Msgf("token %s Nil -> NonExist", id)
			return NonExistToken(s, id)
		}
		log.Debug().Msgf("token %s Nil -> Exist", id)
		return ExistToken(s, id)
	}
}

// NonExistToken represents token that does not exist
func NonExistToken(s *System, id string) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {

		switch msg := context.Data.(type) {

		case SynchornizationDone:
			log.Debug().Msgf("token %s (NonExist CreateToken)", id)
			return NonExistToken(s, id)

		case CreateToken:
			err := persistence.CreateToken(s.EncryptedStorage, id, msg.Value)
			if err != nil {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (NonExist CreateToken) Error %s", id, err)
				return NonExistToken(s, id)
			}

			s.SendMessage(RespCreateToken, context.Sender, context.Receiver)
			log.Info().Msgf("New Token %s Created", id)
			log.Debug().Msgf("token %s (NonExist CreateToken) OK", id)
			s.Metrics.TokenCreated()

			return ExistToken(s, id)

		case DeleteToken:
			s.SendMessage(RespTokenDoesNotExist, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist DeleteToken) Error", id)
			return NonExistToken(s, id)

		case SynchronizeToken:
			s.SendMessage(RespTokenDoesNotExist, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist SynchronizeToken) Error", id)
			return NonExistToken(s, id)

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist Unknown Message) Error", id)
			return NonExistToken(s, id)
		}

	}
}

// ExistToken represents account that does exist
func ExistToken(s *System, id string) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {

		switch context.Data.(type) {

		case SynchornizationDone:
			log.Debug().Msgf("token %s (Synchronizing CreateToken)", id)
			return ExistToken(s, id)

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Exist CreateToken) Error", id)
			return ExistToken(s, id)

		case SynchronizeToken:
			log.Debug().Msgf("token %s (Exist SynchronizeToken)", id)
			log.Info().Msgf("Synchronizing %s", id)
			s.SendMessage(RespSynchronizeToken, context.Sender, context.Receiver)

			go func() {
				log.Debug().Msgf("token %s Importing statements Start", id)

				defer func() {
					log.Debug().Msgf("token %s Importing statements End", id)
					context.Self.Tell(SynchornizationDone{}, context.Receiver, context.Receiver)
				}()

				token, err := persistence.LoadToken(s.EncryptedStorage, id)
				if err != nil {
					// TODO log
					return
				}

				workflow := integration.NewWorkflow(
					token,
					s.Tenant,
					s.FioGateway,
					s.VaultGateway,
					s.LedgerGateway,
					s.EncryptedStorage,
					s.PlaintextStorage,
					s.Metrics,
				)

				//workflow.SynchronizeStatements()

				workflow.DownloadStatements()

				/*
					workflow.DownloadStatements()
					workflow.CreateAccounts()
					workflow.CreateTransactions()
				*/
			}()

			return SynchronizingToken(s, id)

		case DeleteToken:
			if !persistence.DeleteToken(s.EncryptedStorage, id) {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (Exist DeleteToken) Error", id)
				return ExistToken(s, id)
			}

			log.Info().Msgf("Token %s Deleted", id)
			log.Debug().Msgf("token %s (Exist DeleteToken) OK", id)
			s.Metrics.TokenDeleted()
			s.SendMessage(RespDeleteToken, context.Sender, context.Receiver)
			return NonExistToken(s, id)

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Exist Unknown Message) Error", id)
			return ExistToken(s, id)

		}

	}
}

// SynchronizingToken represents account that is currently synchronizing
func SynchronizingToken(s *System, id string) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {

		switch context.Data.(type) {

		case SynchornizationDone:
			log.Debug().Msgf("token %s (Synchronizing SynchornizationDone)", id)
			return ExistToken(s, id)

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Synchronizing CreateToken) Error", id)
			return SynchronizingToken(s, id)

		case SynchronizeToken:
			s.SendMessage(RespSynchronizeToken, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Synchronizing SynchronizeToken)", id)
			return SynchronizingToken(s, id)

		case DeleteToken:
			if !persistence.DeleteToken(s.EncryptedStorage, id) {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (Synchronizing DeleteToken) Error", id)
				return SynchronizingToken(s, id)
			}
			log.Info().Msgf("Token %s Deleted", id)
			log.Debug().Msgf("token %s (Synchronizing DeleteToken) OK", id)
			s.Metrics.TokenDeleted()
			s.SendMessage(RespDeleteToken, context.Sender, context.Receiver)
			return NonExistToken(s, id)

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Warn().Msgf("token %s (Synchronizing Unknown Message) Error", id)
			return SynchronizingToken(s, id)

		}

	}
}
