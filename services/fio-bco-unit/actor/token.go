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

package actor

import (
	"github.com/jancajthaml-openbank/fio-bco-unit/daemon"
	"github.com/jancajthaml-openbank/fio-bco-unit/model"
	"github.com/jancajthaml-openbank/fio-bco-unit/persistence"

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

func TokenManagement(s *daemon.ActorSystem) func(interface{}, system.Context) {
	return func(t_state interface{}, context system.Context) {

		switch msg := context.Data.(type) {

		case model.CreateToken:
			log.Debug("token ~ (CreateToken)")

			if !persistence.CreateToken(s.Storage, msg.Value) {
				s.SendRemote(context.Sender.Region, FatalErrorMessage(context.Receiver.Name, context.Sender.Name))
				log.Debug("token ~ (CreateToken) Error")
				return
			}

			log.Infof("Token %s Created", msg.Value)

			s.Metrics.TokenCreated()

			s.SendRemote(context.Sender.Region, TokenCreatedMessage(context.Receiver.Name, context.Sender.Name))

		case model.DeleteToken:

			if !persistence.DeleteToken(s.Storage, msg.Value) {
				s.SendRemote(context.Sender.Region, FatalErrorMessage(context.Receiver.Name, context.Sender.Name))
				log.Debug("token ~ (DeleteToken) Error")
				return
			}

			log.Infof("Token %s Deleted", msg.Value)

			s.Metrics.TokenDeleted()

			s.SendRemote(context.Sender.Region, TokenDeletedMessage(context.Receiver.Name, context.Sender.Name))

		default:
			s.SendRemote(context.Sender.Region, FatalErrorMessage(context.Receiver.Name, context.Sender.Name))
			log.Debug("token ~ (Unknown Message) Error")
		}

		return
	}
}

func NewTokenSignletonActor() *system.Envelope {
	return system.NewEnvelope("token", nil)
}
