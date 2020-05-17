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

package persistence

import (
	"time"

	localfs "github.com/jancajthaml-openbank/local-fs"

	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/utils"
)

// LoadTokens rehydrates token entity state from storage
func LoadTokens(storage *localfs.EncryptedStorage) ([]model.Token, error) {
	path := utils.TokensPath()
	ok, err := storage.Exists(path)
	if err != nil || !ok {
		return make([]model.Token, 0), nil
	}
	tokens, err := storage.ListDirectory(path, true)
	if err != nil {
		return nil, err
	}
	result := make([]model.Token, len(tokens))
	for i, id := range tokens {
		token := model.Token{
			ID: id,
		}
		if HydrateToken(storage, &token) != nil {
			result[i] = token
		}
	}
	return result, nil
}

// LoadToken rehydrates token entity state from storage
func LoadToken(storage *localfs.EncryptedStorage, id string) *model.Token {
	result := new(model.Token)
	result.ID = id
	return HydrateToken(storage, result)
}

// CreateToken persist token entity state to storage
func CreateToken(storage *localfs.EncryptedStorage, id string, value string) *model.Token {
	return PersistToken(storage, &model.Token{
		ID:           id,
		Value:        value,
		CreatedAt:    time.Now().UTC(),
		LastSyncedID: 0,
	})
}

// DeleteToken deletes existing token entity
func DeleteToken(storage *localfs.EncryptedStorage, id string) bool {
	path := utils.TokenPath(id)
	return storage.DeleteFile(path) == nil
}

// PersistToken persist new token entity to storage
func PersistToken(storage *localfs.EncryptedStorage, entity *model.Token) *model.Token {
	if entity == nil {
		return nil
	}
	path := utils.TokenPath(entity.ID)
	data, err := entity.Serialise()
	if err != nil {
		return nil
	}
	if storage.WriteFileExclusive(path, data) != nil {
		return nil
	}
	return entity
}

// HydrateToken hydrate existing token from storage
func HydrateToken(storage *localfs.EncryptedStorage, entity *model.Token) *model.Token {
	if entity == nil {
		return nil
	}
	path := utils.TokenPath(entity.ID)
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

// UpdateToken updates data of existing token to storage
func UpdateToken(storage *localfs.EncryptedStorage, entity *model.Token) bool {
	if entity == nil {
		return false
	}
	path := utils.TokenPath(entity.ID)
	// FIXME check nil
	data, err := entity.Serialise()
	if err != nil {
		return false
	}
	return storage.WriteFile(path, data) == nil
}
