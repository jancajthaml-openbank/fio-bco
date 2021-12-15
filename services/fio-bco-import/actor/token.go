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
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"

	system "github.com/jancajthaml-openbank/actor-system"
)

// NilToken represents token that is neither existing neither non existing
func NilToken(s *System, state model.Token) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {
		context.Self.Tell(context.Data, context.Receiver, context.Sender)
		tokenHydration := persistence.LoadToken(s.EncryptedStorage, state.ID)
		if tokenHydration == nil {
			log.Debug().Msgf("token %s Nil -> NonExist", state.ID)
			return NonExistToken(s, state)
		} else {
			log.Debug().Msgf("token %s Nil -> Exist", state.ID)
			return ExistToken(s, *tokenHydration)
		}
	}
}

// NonExistToken represents token that does not exist
func NonExistToken(s *System, state model.Token) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {

		switch msg := context.Data.(type) {

		case ProbeMessage:
			return NonExistToken(s, state)

		case CreateToken:
			tokenResult := persistence.CreateToken(s.EncryptedStorage, state.ID, msg.Value)

			if tokenResult == nil {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (NonExist CreateToken) Error", state.ID)
				return NonExistToken(s, state)
			}

			s.SendMessage(RespCreateToken, context.Sender, context.Receiver)
			log.Info().Msgf("New Token %s Created", state.ID)
			log.Debug().Msgf("token %s (NonExist CreateToken) OK", state.ID)
			s.Metrics.TokenCreated()

			return ExistToken(s, *tokenResult)

		case DeleteToken:
			s.SendMessage(RespTokenDoesNotExist, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist DeleteToken) Error", state.ID)
			return NonExistToken(s, state)

		case SynchronizeToken:
			s.SendMessage(RespTokenDoesNotExist, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist SynchronizeToken) Error", state.ID)
			return NonExistToken(s, state)

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (NonExist Unknown Message) Error", state.ID)
			return NonExistToken(s, state)
		}

	}
}

// ExistToken represents account that does exist
func ExistToken(s *System, state model.Token) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {

		switch context.Data.(type) {

		case ProbeMessage:
			return ExistToken(s, state)

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Exist CreateToken) Error", state.ID)
			return ExistToken(s, state)

		case SynchronizeToken:
			log.Debug().Msgf("token %s (Exist SynchronizeToken)", state.ID)
			log.Info().Msgf("Synchronizing %s", state.ID)
			s.SendMessage(RespSynchronizeToken, context.Sender, context.Receiver)

			go func() {
				log.Debug().Msgf("token %s Importing statements Start", state.ID)

				defer func() {
					log.Debug().Msgf("token %s Importing statements End", state.ID)
					context.Self.Tell(ProbeMessage{}, context.Receiver, context.Receiver)
				}()

				workflow := integration.NewWorkflow(
					&state,	// FIXME copy
					s.Tenant,
					s.FioGateway,
					s.VaultGateway,
					s.LedgerGateway,
					s.EncryptedStorage,
					s.PlaintextStorage,
					s.Metrics,
				)

				workflow.SynchronizeStatements()
			}()

			return SynchronizingToken(s, state)

		case DeleteToken:
			if !persistence.DeleteToken(s.EncryptedStorage, state.ID) {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (Exist DeleteToken) Error", state.ID)
				return ExistToken(s, state)
			}

			log.Info().Msgf("Token %s Deleted", state.ID)
			log.Debug().Msgf("token %s (Exist DeleteToken) OK", state.ID)
			s.Metrics.TokenDeleted()
			s.SendMessage(RespDeleteToken, context.Sender, context.Receiver)
			return NonExistToken(s, model.NewToken(state.ID))

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Exist Unknown Message) Error", state.ID)
			return ExistToken(s, state)

		}

	}
}

// SynchronizingToken represents account that is currently synchronizing
func SynchronizingToken(s *System, state model.Token) system.ReceiverFunction {
	return func(context system.Context) system.ReceiverFunction {

		switch context.Data.(type) {

		case ProbeMessage:
			return SynchronizingToken(s, state)

		case CreateToken:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Synchronizing CreateToken) Error", state.ID)
			return SynchronizingToken(s, state)

		case SynchronizeToken:
			s.SendMessage(RespSynchronizeToken, context.Sender, context.Receiver)
			log.Debug().Msgf("token %s (Synchronizing SynchronizeToken)", state.ID)
			return SynchronizingToken(s, state)

		case DeleteToken:
			if !persistence.DeleteToken(s.EncryptedStorage, state.ID) {
				s.SendMessage(FatalError, context.Sender, context.Receiver)
				log.Debug().Msgf("token %s (Synchronizing DeleteToken) Error", state.ID)
				return SynchronizingToken(s, state)
			}
			log.Info().Msgf("Token %s Deleted", state.ID)
			log.Debug().Msgf("token %s (Synchronizing DeleteToken) OK", state.ID)
			s.Metrics.TokenDeleted()
			s.SendMessage(RespDeleteToken, context.Sender, context.Receiver)
			return NonExistToken(s, model.NewToken(state.ID))

		default:
			s.SendMessage(FatalError, context.Sender, context.Receiver)
			log.Warn().Msgf("token %s (Synchronizing Unknown Message) Error", state.ID)
			return SynchronizingToken(s, state)

		}

	}
}
