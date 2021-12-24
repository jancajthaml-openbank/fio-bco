// Copyright (c) 2016-2021, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"fmt"

	localfs "github.com/jancajthaml-openbank/local-fs"

	"github.com/jancajthaml-openbank/fio-bco-import/model"
)

// LoadTokens rehydrates token entity state from storage
func LoadTokens(storage localfs.Storage) ([]model.Token, error) {
	path := "token"
	ok, err := storage.Exists(path)
	if err != nil || !ok {
		return make([]model.Token, 0), nil
	}
	tokens, err := storage.ListDirectory(path, true)
	if err != nil {
		return nil, err
	}
	result := make([]model.Token, 0)
	for _, id := range tokens {
		token := model.Token{
			ID: id,
		}
		if HydrateToken(storage, &token) != nil {
			result = append(result, token)
		}
	}
	return result, nil
}

// LoadToken rehydrates token entity state from storage
func LoadToken(storage localfs.Storage, id string) (*model.Token, error) {
	result := new(model.Token)
	result.ID = id
	err := HydrateToken(storage, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateToken persist token entity state to storage
func CreateToken(storage localfs.Storage, id string, value string) error {
	token := model.NewToken(id)
	token.Value = value
	return PersistToken(storage, &token)
}

// DeleteToken deletes existing token entity
func DeleteToken(storage localfs.Storage, id string) bool {
	path := "token/" + id + "/value"
	return storage.DeleteFile(path) == nil
}

// PersistToken persist new token entity to storage
func PersistToken(storage localfs.Storage, entity *model.Token) error {
	if entity == nil {
		return fmt.Errorf("nil reference")
	}
	data, err := entity.Serialize()
	if err != nil {
		return err
	}
	return storage.WriteFileExclusive("token/"+entity.ID+"/value", data)
}

// HydrateToken hydrate existing token from storage
func HydrateToken(storage localfs.Storage, entity *model.Token) error {
	if entity == nil {
		return fmt.Errorf("nil reference")
	}
	path := "token/" + entity.ID + "/value"
	data, err := storage.ReadFileFully(path)
	if err != nil {
		return err
	}
	return entity.Deserialize(data)
}

// UpdateToken updates data of existing token to storage
func UpdateToken(storage localfs.Storage, entity *model.Token) bool {
	if entity == nil {
		return false
	}
	path := "token/" + entity.ID + "/value"
	data, err := entity.Serialize()
	if err != nil {
		return false
	}
	return storage.WriteFile(path, data) == nil
}
