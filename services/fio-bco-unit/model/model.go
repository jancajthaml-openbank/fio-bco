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

package model

import (
	"bytes"
	"strconv"
	"strings"
)

// Token represents metadata of token entity
type Token struct {
	Value        string
	LastSyncedID int64
}

// ListTokens is inbound request for listing of existing tokens
type ListTokens struct {
}

// CreateToken is inbound request for creation of new token
type CreateToken struct {
	Value string
}

// DeleteToken is inbound request for deletion of token
type DeleteToken struct {
	Value string
}

// GetToken is inbound request for token details
type GetToken struct {
}

// NewToken returns new Token
func NewToken(value string) Token {
	return Token{
		Value:        value,
		LastSyncedID: 0,
	}
}

// Persist serializes Token entity to persistable data
func (entity *Token) Persist() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(strconv.FormatInt(entity.LastSyncedID, 10))

	return buffer.Bytes()
}

// Hydrate deserializes Token entity from persistent data
func (entity *Token) Hydrate(data []byte) {
	lines := strings.Split(string(data), "\n")

	if cast, err := strconv.ParseInt(lines[0], 10, 64); err == nil {
		entity.LastSyncedID = cast
	}
}
