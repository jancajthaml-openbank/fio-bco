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
	"fmt"

	system "github.com/jancajthaml-openbank/actor-system"
)

func parseMessage(msg string, to system.Coordinates) (interface{}, error) {
	start := 0
	end := len(msg)
	parts := make([]string, 2)
	idx := 0
	i := 0
	for i < end && idx < 2 {
		if msg[i] == ' ' {
			if !(start == i && msg[start] == ' ') {
				parts[idx] = msg[start:i]
				idx++
			}
			start = i + 1
		}
		i++
	}
	if idx < 2 && msg[start] != ' ' && len(msg[start:]) > 0 {
		parts[idx] = msg[start:]
		idx++
	}

	if i != end {
		return nil, fmt.Errorf("message too large")
	}

	switch parts[0] {

	case ReqSynchronizeToken:
		return SynchronizeToken{}, nil

	case ReqCreateToken:
		if idx == 2 {
			return CreateToken{
				ID:    to.Name,
				Value: parts[1],
			}, nil
		}
		return nil, fmt.Errorf("invalid message %s", msg)

	case ReqDeleteToken:
		return DeleteToken{
			ID: to.Name,
		}, nil

	default:
		return nil, fmt.Errorf("unknown message %s", msg)
	}
}

// ProcessMessage processing of remote message to this fio-bco
func ProcessMessage(s *System) system.ProcessMessage {
	return func(msg string, to system.Coordinates, from system.Coordinates) {
		message, err := parseMessage(msg, to)
		if err != nil {
			if from != to && to.Name != "" {
				log.Warn().Err(err).Msgf("Failed to parse message [remote %v -> local %v]", from, to)
				s.SendMessage(FatalError, from, to)
			}
			return
		}
		ref, err := s.ActorOf(to.Name)
		if err != nil {
			ref, err = spawnTokenActor(s, to.Name)
		}
		if err != nil {
			if from != to && to.Name != "" {
				log.Warn().Err(err).Msgf("Deadletter [remote %v -> local %v] %s", from, to, msg)
				s.SendMessage(FatalError, to, from)
			}
			return
		}
		ref.Tell(message, to, from)
	}
}

func spawnTokenActor(s *System, id string) (*system.Actor, error) {
	envelope := system.NewActor(id, NilToken(s, id))

	err := s.RegisterActor(envelope)
	if err != nil {
		log.Warn().Msgf("Unable to register %s actor", id)
		return nil, err
	}

	log.Debug().Msgf("Actor %s registered", id)
	return envelope, nil
}
