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
	"time"

	"github.com/jancajthaml-openbank/fio-bco-rest/daemon"
	"github.com/jancajthaml-openbank/fio-bco-rest/model"

	"github.com/rs/xid"

	system "github.com/jancajthaml-openbank/actor-system"
	log "github.com/sirupsen/logrus"
)

// CreateToken creates new token for target tenant
func CreateToken(s *daemon.ActorSystem, tenant string, token string) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("CreateToken recovered in %v", r)
			result = nil
		}
	}()

	// FIXME check if token or tenant is empty string and return error if so

	ch := make(chan interface{})
	defer close(ch)

	envelope := system.NewEnvelope("relay/"+xid.New().String(), nil)
	defer s.UnregisterActor(envelope.Name)

	err := s.RegisterActor(envelope, func(state interface{}, context system.Context) {
		switch msg := context.Data.(type) {
		case model.TokenCreated:
			ch <- &msg
		default:
			ch <- nil
		}
	})
	if err != nil {
		return
	}

	s.SendRemote("FioUnit/"+tenant, CreateTokenMessage(envelope.Name, token))

	select {

	case result = <-ch:
		return

	case <-time.After(time.Second):
		result = new(model.ReplyTimeout)
		return
	}
}

// DeleteToken deletes existing token for target tenant
func DeleteToken(s *daemon.ActorSystem, tenant string, token string) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("DeleteToken recovered in %v", r)
			result = nil
		}
	}()

	// FIXME check if token or tenant is empty string and return error if so

	ch := make(chan interface{})
	defer close(ch)

	envelope := system.NewEnvelope("relay/"+xid.New().String(), nil)
	defer s.UnregisterActor(envelope.Name)

	err := s.RegisterActor(envelope, func(state interface{}, context system.Context) {
		switch msg := context.Data.(type) {
		case model.TokenDeleted:
			ch <- &msg
		default:
			ch <- nil
		}
	})
	if err != nil {
		return
	}

	s.SendRemote("FioUnit/"+tenant, DeleteTokenMessage(envelope.Name, token))

	select {

	case result = <-ch:
		return

	case <-time.After(time.Second):
		result = new(model.ReplyTimeout)
		return
	}
}
