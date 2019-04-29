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

	"github.com/jancajthaml-openbank/fio-bco-import/daemon"
	"github.com/jancajthaml-openbank/fio-bco-import/model"

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

var nilCoordinates = system.Coordinates{}

func asEnvelopes(s *daemon.ActorSystem, parts []string) (system.Coordinates, system.Coordinates, string, error) {
	if len(parts) < 4 {
		return nilCoordinates, nilCoordinates, "", fmt.Errorf("invalid message received %+v", parts)
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

	return from, to, payload, nil
}

func spawnTokenActor(s *daemon.ActorSystem, id string) (*system.Envelope, error) {
	envelope := system.NewEnvelope(id, model.NewToken(id))

	err := s.RegisterActor(envelope, NilToken(s))
	if err != nil {
		log.Warnf("%s ~ Spawning Token Error unable to register", id)
		return nil, err
	}

	log.Debugf("%s ~ Token Spawned", id)
	return envelope, nil
}

// ProcessRemoteMessage processing of remote message to this bondster-bco
func ProcessRemoteMessage(s *daemon.ActorSystem) system.ProcessRemoteMessage {
	return func(parts []string) {
		from, to, payload, err := asEnvelopes(s, parts)
		if err != nil {
			log.Warn(err.Error())
			return
		}

		defer func() {
			if r := recover(); r != nil {
				log.Errorf("procesRemoteMessage recovered in [remote %v -> local %v] : %+v", from, to, r)
				s.SendRemote(from.Region, FatalErrorMessage(to.Name, from.Name))
			}
		}()

		ref, err := s.ActorOf(to.Name)
		if err != nil {
			ref, err = spawnTokenActor(s, to.Name)
		}

		if err != nil {
			log.Warnf("Actor not found [remote %v -> local %v]", from, to)
			s.SendRemote(from.Region, FatalErrorMessage(to.Name, from.Name))
			return
		}

		var message interface{}

		switch payload {

		case ReqCreateToken:
			if len(parts) == 5 {
				message = model.CreateToken{
					ID:    parts[2],
					Value: parts[4],
				}
			} else {
				message = nil
			}

		case ReqDeleteToken:
			message = model.DeleteToken{
				ID: parts[2],
			}

		default:
			message = nil
		}

		if message == nil {
			log.Warnf("Deserialization of unsuported message [remote %v -> local %v] : %+v", from, to, parts)
			s.SendRemote(from.Region, FatalErrorMessage(to.Name, from.Name))
			return
		}

		ref.Tell(message, from)
	}
}

// ProcessLocalMessage processing of local message to this bco
func ProcessLocalMessage(s *daemon.ActorSystem) system.ProcessLocalMessage {
	return func(message interface{}, to system.Coordinates, from system.Coordinates) {
		if to.Region != "" && to.Region != s.Name {
			log.Warnf("Invalid region received [local %s -> local %s]", from, to)
			return
		}
		ref, err := s.ActorOf(to.Name)
		if err != nil {
			ref, err = spawnTokenActor(s, to.Name)
		}

		if err != nil {
			log.Warnf("Actor not found [local %s]", to)
			return
		}
		ref.Tell(message, from)
	}
}
