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

package api

import (
	"io/ioutil"

	"github.com/jancajthaml-openbank/fio-bco-rest/actor"
	"github.com/jancajthaml-openbank/fio-bco-rest/config"
	"github.com/jancajthaml-openbank/fio-bco-rest/daemon"
	"github.com/jancajthaml-openbank/fio-bco-rest/model"
	"github.com/jancajthaml-openbank/fio-bco-rest/persistence"
	"github.com/jancajthaml-openbank/fio-bco-rest/utils"

	"net/http"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

var emptyJSONObject = []byte("{}")
var emptyJSONArray = []byte("[]")

// HealtCheck returns 200 OK
func HealtCheck(w http.ResponseWriter, r *http.Request) {
	log.Debug("HealtCheck request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(emptyJSONObject)
}

// CreateTokenPartial returns http handler for creating new token
func CreateTokenPartial(system *daemon.ActorSystem) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Debug("CreateToken request")

		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		tenant := vars["tenant_id"]

		// Read body
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(emptyJSONObject)
			return
		}

		var req model.Token
		err = utils.JSON.Unmarshal(b, &req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(emptyJSONObject)
			return
		}

		result := actor.CreateToken(system, req.Value, tenant)

		switch result.(type) {

		case *model.TokenCreated:
			log.Debug("Smiley ok here")

			w.WriteHeader(http.StatusOK)
			w.Write(emptyJSONObject)
			return

		case *model.ReplyTimeout:
			log.Debug("Sad timeout here")

			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write(emptyJSONObject)
			return

		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(emptyJSONObject)
			return

		}
	}
}

// DeleteTokenPartial returns http handler for deleting existing token
func DeleteTokenPartial(system *daemon.ActorSystem) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Debug("DeleteToken request")

		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		tenant := vars["tenant_id"]
		token := vars["token_id"]

		result := actor.DeleteToken(system, token, tenant)

		switch result.(type) {

		case *model.TokenDeleted:
			w.WriteHeader(http.StatusOK)
			w.Write(emptyJSONObject)
			return

		case *model.ReplyTimeout:
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write(emptyJSONObject)
			return

		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(emptyJSONObject)
			return

		}
	}
}

// GetTokensPartial returns http handler for getting tokens
func GetTokensPartial(cfg config.Configuration, system *daemon.ActorSystem) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Debug("GetTokens request")

		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		tokens, err := persistence.LoadTokens(cfg.RootStorage, vars["tenant_id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(emptyJSONArray)
			return
		}

		resp, err := utils.JSON.Marshal(tokens)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(emptyJSONArray)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)

		return
	}
}
