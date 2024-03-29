// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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

func parseMessage(msg string) (interface{}, error) {
	end := len(msg)
	i := 0
	for i < end && msg[i] != ' ' {
		i++
	}

	switch msg[0:i] {

	case FatalError:
		return FatalError, nil

	case RespCreateToken:
		return new(TokenCreated), nil

	case RespDeleteToken:
		return new(TokenDeleted), nil

	case RespSynchronizeToken:
		return new(TokenSynchonizeAccepted), nil

	case RespTokenDoesNotExist:
		return new(TokenMissing), nil

	default:
		return nil, fmt.Errorf("unknown message %s", msg)
	}
}

// ProcessMessage processing of remote message
func ProcessMessage(s *System) system.ProcessMessage {
	return func(msg string, to system.Coordinates, from system.Coordinates) {
		ref, err := s.ActorOf(to.Name)
		if err != nil {
			log.Warn().Err(err).Msgf("Deadletter [remote %v -> local %v] %s", from, to, msg)
			return
		}
		message, err := parseMessage(msg)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to parse message [remote %v -> local %v]", from, to)
		}
		ref.Tell(message, to, from)
		return
	}
}
