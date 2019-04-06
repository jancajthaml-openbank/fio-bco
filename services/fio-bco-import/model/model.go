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

package model

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Token represents metadata of token entity
type Token struct {
	ID           string
	Value        string
	LastSyncedID int64
}

// ListTokens is inbound request for listing of existing tokens
type ListTokens struct {
}

// CreateToken is inbound request for creation of new token
type CreateToken struct {
	ID    string
	Value string
}

// DeleteToken is inbound request for deletion of token
type DeleteToken struct {
	ID string
}

// GetToken is inbound request for token details
type GetToken struct {
}

// NewToken returns new Token
func NewToken(id string, value string) Token {
	return Token{
		ID:           id,
		Value:        value,
		LastSyncedID: 0,
	}
}

// Serialise Token entity to persistable data
func (entity *Token) Serialise() ([]byte, error) {
	if entity == nil {
		return nil, fmt.Errorf("called Token.Serialise over nil")
	}
	var buffer bytes.Buffer
	buffer.WriteString(entity.Value)
	buffer.WriteString("\n")
	buffer.WriteString(strconv.FormatInt(entity.LastSyncedID, 10))
	return buffer.Bytes(), nil
}

// Deserialise Token entity from persistent data
func (entity *Token) Deserialise(data []byte) error {
	if entity == nil {
		return fmt.Errorf("called Token.Deserialise over nil")
	}

	// FIXME more optimal split
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("malformed data")
	}

	entity.Value = lines[0]

	if cast, err := strconv.ParseInt(lines[1], 10, 64); err == nil {
		entity.LastSyncedID = cast
	}

	return nil
}
