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

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

// ProcessLocalMessage processing of local message to this fio-bco
func ProcessLocalMessage(s *daemon.ActorSystem) system.ProcessLocalMessage {
	return func(message interface{}, to system.Coordinates, from system.Coordinates) {
		if to.Region != "" && to.Region != s.Name {
			log.Warnf("Invalid region received [local %s -> local %s]", from, to)
			return
		}

		log.Debugf("Inherited Actor System received local message %+v", message)
	}
}

// ProcessRemoteMessage processing of remote message to this fio-bco
func ProcessRemoteMessage(s *daemon.ActorSystem) system.ProcessRemoteMessage {
	return func(parts []string) {
		if len(parts) < 4 {
			log.Warnf("invalid message received %+v", parts)
			return
		}

		region, receiver, sender, payload := parts[0], parts[1], parts[2], parts[3]

		from := system.Coordinates{
			Name:   sender,
			Region: region,
		}

		to := system.Coordinates{
			Name:   receiver,
			Region: s.Name,
		}

		defer func() {
			if r := recover(); r != nil {
				log.Errorf("procesRemoteMessage recovered in [remote %v -> local %v] : %+v", from, to, r)
			}
		}()

		ref, err := s.ActorOf(to.Name)
		if err != nil {
			log.Warnf("Actor not found [remote %v -> local %v]", from, to)
			return
		}

		var message interface{}

		switch payload {

		case ReqCreateToken:
			if len(parts) == 5 {
				message = model.CreateToken{
					Value: parts[4],
				}
			} else {
				message = nil
			}

		case ReqDeleteToken:
			if len(parts) == 5 {
				message = model.DeleteToken{
					Value: parts[4],
				}
			} else {
				message = nil
			}

		default:
			message = nil
		}

		if message == nil {
			log.Warnf("Deserialization of unsuported message [remote %v -> local %v] : %+v", from, to, parts)
			s.SendRemote(region, FatalErrorMessage(to.Name, from.Name))
			return
		}

		ref.Tell(message, from)
	}
}

// SpawnTokenSignletonActor spawns actor for token CRUD operations
func SpawnTokenSignletonActor(s *daemon.ActorSystem) (*system.Envelope, error) {
	envelope := NewTokenSignletonActor()

	err := s.RegisterActor(envelope, TokenManagement(s))
	if err != nil {
		log.Warn("token ~ Spawning Actor Error unable to register")
		return nil, err
	}

	log.Debug("token ~ Actor Spawned")
	return envelope, nil
}
