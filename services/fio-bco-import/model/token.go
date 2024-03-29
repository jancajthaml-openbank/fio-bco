// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"time"
)

// Token represents metadata of token entity
type Token struct {
	ID           string
	Value        string
	CreatedAt    time.Time
	LastSyncedID int64
}

// NewToken returns new Token
func NewToken(id string) Token {
	return Token{
		ID:           id,
		Value:        "",
		CreatedAt:    time.Now().UTC(),
		LastSyncedID: 0,
	}
}

// Serialize Token entity to persistable data
func (entity *Token) Serialize() ([]byte, error) {
	if entity == nil {
		return nil, fmt.Errorf("called Token.Serialize over nil")
	}
	var buffer bytes.Buffer
	buffer.WriteString(entity.CreatedAt.Format(time.RFC3339))
	buffer.WriteString("\n")
	buffer.WriteString(entity.Value)
	buffer.WriteString("\n")
	buffer.WriteString(strconv.FormatInt(entity.LastSyncedID, 10))
	return buffer.Bytes(), nil
}

// Deserialize Token entity from persistent data
func (entity *Token) Deserialize(data []byte) error {
	if entity == nil {
		return fmt.Errorf("called Token.Deserialize over nil")
	}

	// FIXME more optimal split
	lines := strings.Split(string(data), "\n")
	if len(lines) < 3 {
		return fmt.Errorf("malformed data")
	}

	if cast, err := time.Parse(time.RFC3339, lines[0]); err == nil {
		entity.CreatedAt = cast
	}

	entity.Value = lines[1]

	if cast, err := strconv.ParseInt(lines[2], 10, 64); err == nil {
		entity.LastSyncedID = cast
	}

	return nil
}
