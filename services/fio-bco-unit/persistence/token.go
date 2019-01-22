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

	"github.com/jancajthaml-openbank/fio-bco-unit/model"
	"github.com/jancajthaml-openbank/fio-bco-unit/utils"
)

// LoadTokens rehydrates token entity state from storage
func LoadTokens(storage *localfs.Storage) ([]model.Token, error) {
	path := utils.TokensPath()
	tokens, err := storage.ListDirectory(path, true)
	if err != nil {
		return nil, err
	}
	result := make([]model.Token, len(tokens))
	for i, value := range tokens {
		token := model.Token{
			Value: value,
		}
		if HydrateToken(storage, &token) != nil {
			result[i] = token
		}
	}
	return result, nil
}

// CreateToken creates and persist new token entity
func CreateToken(storage *localfs.Storage, value string) bool {
	token := model.NewToken(value)
	if PersistToken(storage, &token) == nil {
		return false
	}
	return true
}

// DeleteToken deletes existing token entity
func DeleteToken(storage *localfs.Storage, value string) bool {
	path := utils.TokenPath(value)
	return storage.DeleteFile(path) == nil
}

// PersistToken persist new token entity to storage
func PersistToken(storage *localfs.Storage, entity *model.Token) *model.Token {
	if entity == nil {
		return nil
	}
	path := utils.TokenPath(entity.Value)
	if storage.TouchFile(path) != nil {
		return nil
	}
	return entity
}

// HydrateToken hydrate existing token from storage
func HydrateToken(storage *localfs.Storage, entity *model.Token) *model.Token {
	if entity == nil {
		return nil
	}
	path := utils.TokenPath(entity.Value)
	data, err := storage.ReadFileFully(path)
	if err != nil {
		return nil
	}
	entity.Hydrate(data)
	return entity
}

// UpdateToken updates data of existing token to storage
func UpdateToken(storage *localfs.Storage, entity *model.Token) bool {
	if entity == nil {
		return false
	}
	path := utils.TokenPath(entity.Value)
	data := entity.Persist()
	return storage.UpdateFile(path, data) == nil
}
