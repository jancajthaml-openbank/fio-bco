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
	"time"

	"github.com/rs/xid"

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

// CreateToken creates new token for target tenant
func CreateToken(sys *ActorSystem, tenant string, token Token) (result interface{}) {
	sys.Metrics.TimeCreateToken(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("CreateToken recovered in %v", r)
				result = nil
			}
		}()

		ch := make(chan interface{})
		defer close(ch)

		envelope := system.NewEnvelope("relay/"+xid.New().String(), nil)
		defer sys.UnregisterActor(envelope.Name)

		sys.RegisterActor(envelope, func(state interface{}, context system.Context) {
			switch msg := context.Data.(type) {
			case *TokenCreated:
				ch <- msg
			default:
				ch <- nil
			}
		})

		sys.SendRemote(CreateTokenMessage(tenant, envelope.Name, token))

		select {

		case result = <-ch:
			log.Infof("Token %s/%s created", tenant, token.ID)
			return

		case <-time.After(time.Second):
			result = new(ReplyTimeout)
			return
		}
	})
	return
}

// DeleteToken deletes existing token for target tenant
func DeleteToken(sys *ActorSystem, tenant string, token string) (result interface{}) {
	sys.Metrics.TimeDeleteToken(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("DeleteToken recovered in %v", r)
				result = nil
			}
		}()

		ch := make(chan interface{})
		defer close(ch)

		envelope := system.NewEnvelope("relay/"+xid.New().String(), nil)
		defer sys.UnregisterActor(envelope.Name)

		sys.RegisterActor(envelope, func(state interface{}, context system.Context) {
			switch msg := context.Data.(type) {
			case *TokenDeleted:
				log.Infof("Token %s/%s deleted", tenant, token)
				ch <- msg
			default:
				ch <- nil
			}
		})

		sys.SendRemote(DeleteTokenMessage(tenant, envelope.Name, token))

		select {

		case result = <-ch:
			return

		case <-time.After(time.Second):
			result = new(ReplyTimeout)
			return
		}
	})
	return
}
