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

package persistence

import (
	localfs "github.com/jancajthaml-openbank/local-fs"

	"github.com/jancajthaml-openbank/fio-bco-rest/actor"
	"github.com/jancajthaml-openbank/fio-bco-rest/utils"
)

// LoadTokens rehydrates token entity state from storage
func LoadTokens(storage *localfs.EncryptedStorage, tenant string) ([]actor.Token, error) {
	path := utils.TokensPath(tenant)
	ok, err := storage.Exists(path)
	if err != nil || !ok {
		return make([]actor.Token, 0), nil
	}
	tokens, err := storage.ListDirectory(path, true)
	if err != nil {
		return nil, err
	}
	var result = make([]actor.Token, 0)
	for _, id := range tokens {
		token := actor.Token{
			ID: id,
		}
		if HydrateToken(storage, tenant, &token) != nil {
			result = append(result, token)
		}
	}
	return result, nil
}

// HydrateToken hydrate existing token from storage
func HydrateToken(storage *localfs.EncryptedStorage, tenant string, entity *actor.Token) *actor.Token {
	if entity == nil {
		return nil
	}
	path := utils.TokenPath(tenant, entity.ID)
	data, err := storage.ReadFileFully(path)
	if err != nil {
		return nil
	}
	err = entity.Deserialise(data)
	if err != nil {
		return nil
	}
	return entity
}
