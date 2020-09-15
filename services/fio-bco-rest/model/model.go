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

package model

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"encoding/json"
	"github.com/rs/xid"
)

// Token represents metadata of token entity
type Token struct {
	ID        string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	Value     string    `json:"value"`
}

// MarshalJSON serializes Token as json
func (entity Token) MarshalJSON() ([]byte, error) {
	return []byte("{\"id\":\"" + entity.ID + "\",\"createdAt\":\"" + entity.CreatedAt.Format(time.RFC3339) + "\"}"), nil
}

// UnmarshalJSON unmarshal json of Token entity
func (entity *Token) UnmarshalJSON(data []byte) error {
	if entity == nil {
		return fmt.Errorf("cannot unmarshal to nil pointer")
	}
	all := struct {
		Value string `json:"value"`
	}{}
	err := json.Unmarshal(data, &all)
	if err != nil {
		return err
	}
	if all.Value == "" {
		return fmt.Errorf("missing attribute \"value\"")
	}
	entity.Value = all.Value

	noise := make([]byte, 10)
	rand.Read(noise)
	entity.ID = hex.EncodeToString(noise) + xid.New().String()

	return nil
}

// Deserialize Token entity from persistent data
func (entity *Token) Deserialize(data []byte) error {
	if entity == nil {
		return fmt.Errorf("called Token.Deserialize over nil")
	}

	// FIXME more optimal split
	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		return fmt.Errorf("malformed data")
	}

	if cast, err := time.Parse(time.RFC3339, lines[0]); err == nil {
		entity.CreatedAt = cast
	}

	return nil
}
