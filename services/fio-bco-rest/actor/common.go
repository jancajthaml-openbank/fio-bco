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
	"github.com/jancajthaml-openbank/fio-bco-rest/daemon"
	"github.com/jancajthaml-openbank/fio-bco-rest/model"

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

// ProcessRemoteMessage processing of remote message to this wall
func ProcessRemoteMessage(s *daemon.ActorSystem) system.ProcessRemoteMessage {
	return func(parts []string) {
		if len(parts) < 4 {
			log.Warnf("invalid message received %+v", parts)
			return
		}

		region, receiver, sender, payload := parts[0], parts[1], parts[2], parts[3]

		// FIXME receiver and sender are swapped
		from := system.Coordinates{
			Name:   receiver, //sender,
			Region: region,
		}

		to := system.Coordinates{
			Name:   sender, //receiver,
			Region: s.Name,
		}

		defer func() {
			if r := recover(); r != nil {
				log.Errorf("procesRemoteMessage recovered in [remote %v -> local %v] : %+v", from, to, r)
			}
		}()

		ref, err := s.ActorOf(to.Name)
		if err != nil {
			// FIXME forward into deadletter receiver and finish whatever has started
			log.Warnf("Deadletter received [remote %v -> local %v] : %+v", from, to, parts[3:])
			return
		}

		var message interface{}

		switch payload {

		case FatalError:
			message = FatalError

		case RespCreateToken:
			message = model.TokenCreated{}

		case RespDeleteToken:
			message = model.TokenDeleted{}

		default:
			log.Warnf("Deserialization of unsuported message [remote %v -> local %v] : %+v", from, to, parts)
			message = FatalError
		}

		ref.Tell(message, from)
		return
	}
}
