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

package persistence

import (
	localfs "github.com/jancajthaml-openbank/local-fs"

	"github.com/jancajthaml-openbank/fio-bco-rest/model"
	"github.com/jancajthaml-openbank/fio-bco-rest/utils"
)

// LoadTokens rehydrates token entity state from storage
func LoadTokens(storage *localfs.Storage, tenant string) ([]model.Token, error) {
	path := utils.TokensPath(tenant)
	ok, err := storage.Exists(path)
	if err != nil || !ok {
		return make([]model.Token, 0), nil
	}
	tokens, err := storage.ListDirectory(path, true)
	if err != nil {
		return nil, err
	}
	result := make([]model.Token, len(tokens))
	for i, value := range tokens {
		token := model.Token{
			Value: value,
		}
		if HydrateToken(storage, tenant, &token) != nil {
			result[i] = token
		}
	}
	return result, nil
}

// HydrateToken hydrate existing token from storage
func HydrateToken(storage *localfs.Storage, tenant string, entity *model.Token) *model.Token {
	if entity == nil {
		return nil
	}
	path := utils.TokenPath(tenant, entity.Value)
	data, err := storage.ReadFileFully(path)
	if err != nil {
		return nil
	}
	entity.Hydrate(data)
	return entity
}
