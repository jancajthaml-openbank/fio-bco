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
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/jancajthaml-openbank/fio-bco-rest/actor"
	"github.com/jancajthaml-openbank/fio-bco-rest/daemon"
	"github.com/jancajthaml-openbank/fio-bco-rest/model"
	"github.com/jancajthaml-openbank/fio-bco-rest/persistence"
	"github.com/jancajthaml-openbank/fio-bco-rest/utils"

	"github.com/gorilla/mux"
	localfs "github.com/jancajthaml-openbank/local-fs"
	"github.com/rs/xid"
)

// TokenPartial returns http handler for single token
func TokenPartial(system *daemon.ActorSystem) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tenant := vars["tenant"]
		token := vars["token"]

		if tenant == "" || token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(emptyJSONObject)
			return
		}

		switch r.Method {

		case "DELETE":
			DeleteToken(system, tenant, token, w, r)
			return

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(emptyJSONObject)
			return

		}
	}
}

// TokensPartial returns http handler for tokens
func TokensPartial(system *daemon.ActorSystem, storage *localfs.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		tenant := vars["tenant"]

		if tenant == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(emptyJSONArray)
			return
		}

		switch r.Method {

		case "GET":
			GetTokens(storage, tenant, w, r)
			return

		case "POST":
			CreateToken(system, tenant, w, r)
			return

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(emptyJSONObject)
			return

		}

	}
}

// CreateToken creates new token
func CreateToken(system *daemon.ActorSystem, tenant string, w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(emptyJSONObject)
		return
	}

	var req = new(model.Token)
	err = utils.JSON.Unmarshal(b, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(emptyJSONObject)
		return
	}

	noise := make([]byte, 10)
	rand.Read(noise)
	req.ID = hex.EncodeToString(noise) + xid.New().String()

	switch actor.CreateToken(system, tenant, *req).(type) {

	case *model.TokenCreated:

		resp, err := utils.JSON.Marshal(req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(emptyJSONArray)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		return

	case *model.ReplyTimeout:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write(emptyJSONObject)
		return

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(emptyJSONObject)
		return

	}
}

// DeleteToken deletes existing token
func DeleteToken(system *daemon.ActorSystem, tenant string, token string, w http.ResponseWriter, r *http.Request) {
	switch actor.DeleteToken(system, tenant, token).(type) {

	case *model.TokenDeleted:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(emptyJSONObject)
		return

	case *model.ReplyTimeout:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write(emptyJSONObject)
		return

	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(emptyJSONObject)
		return

	}
}

// GetTokens retruns list of existing tokens
func GetTokens(storage *localfs.Storage, tenant string, w http.ResponseWriter, r *http.Request) {
	tokens, err := persistence.LoadTokens(storage, tenant)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(emptyJSONArray)
		return
	}

	resp, err := utils.JSON.Marshal(tokens)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(emptyJSONArray)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}