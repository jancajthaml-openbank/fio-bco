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
	"github.com/jancajthaml-openbank/fio-bco-unit/model"
	"github.com/jancajthaml-openbank/fio-bco-unit/utils"
)

// LoadTokens rehydrates token entity state from storage
func LoadTokens(root string) ([]model.Token, error) {
	path := utils.TokensPath(root)
	tokens, err := utils.ListDirectory(path, true)
	if err != nil {
		return nil, err
	}
	result := make([]model.Token, len(tokens))
	for i, value := range tokens {

		token := model.Token{
			Value: value,
		}

		if HydrateToken(root, &token) != nil {
			result[i] = token
		}

	}

	return result, nil
}

// CreateToken creates and persist new token entity
func CreateToken(root, value string) bool {
	token := model.NewToken(value)
	if PersistToken(root, &token) == nil {
		return false
	}
	return true
}

// DeleteToken deletes existing token entity
func DeleteToken(root, value string) bool {
	path := utils.TokenPath(root, value)
	return utils.DeleteFile(path)
}

// PersistToken persist new token entity to storage
func PersistToken(root string, entity *model.Token) *model.Token {
	if entity == nil {
		return nil
	}

	path := utils.TokenPath(root, entity.Value)

	if !utils.TouchFile(path) {
		return nil
	}

	return entity
}

// HydrateToken hydrate existing token from storage
func HydrateToken(root string, entity *model.Token) *model.Token {
	if entity == nil {
		return nil
	}

	path := utils.TokenPath(root, entity.Value)

	data, err := utils.ReadFileFully(path)
	if err != nil {
		return nil
	}

	entity.Hydrate(data)

	return entity
}

// UpdateToken updates data of existing token to storage
func UpdateToken(root string, entity *model.Token) bool {
	if entity == nil {
		return false
	}

	path := utils.TokenPath(root, entity.Value)
	data := entity.Persist()

	return utils.UpdateFile(path, data)
}
