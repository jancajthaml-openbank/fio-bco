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
	"strings"

	"github.com/jancajthaml-openbank/fio-bco-import/model"

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

// ProcessRemoteMessage processing of remote message to this bondster-bco
func ProcessMessage(s *ActorSystem) system.ProcessMessage {
	return func(msg string, to system.Coordinates, from system.Coordinates) {

		defer func() {
			if r := recover(); r != nil {
				log.Errorf("procesRemoteMessage recovered in [remote %v -> local %v] : %+v", from, to, r)
				s.SendMessage(FatalErrorMessage(), to, from)
			}
		}()

		ref, err := s.ActorOf(to.Name)
		if err != nil {
			ref, err = spawnTokenActor(s, to.Name)
		}

		if err != nil {
			log.Warnf("Actor not found [remote %v -> local %v]", from, to)
			s.SendMessage(FatalErrorMessage(), to, from)
			return
		}

		parts := strings.Split(msg, " ")

		var message interface{}

		switch parts[0] {

		case SynchronizeTokens:
			message = model.SynchronizeToken{}

		case ReqCreateToken:
			if len(parts) == 2 {
				message = model.CreateToken{
					ID:    to.Name,
					Value: parts[1],
				}
			} else {
				message = nil
			}

		case ReqDeleteToken:
			message = model.DeleteToken{
				ID: to.Name,
			}

		default:
			message = nil
		}

		if message == nil {
			log.Warnf("Deserialization of unsuported message [remote %v -> local %v] : %+v", from, to, parts)
			s.SendMessage(FatalErrorMessage(), to, from)
			return
		}

		ref.Tell(message, to, from)
	}
}

func spawnTokenActor(s *ActorSystem, id string) (*system.Envelope, error) {
	envelope := system.NewEnvelope(id, model.NewToken(id))

	err := s.RegisterActor(envelope, NilToken(s))
	if err != nil {
		log.Warnf("%s ~ Spawning Token Error unable to register", id)
		return nil, err
	}

	log.Debugf("%s ~ Token Spawned", id)
	return envelope, nil
}
