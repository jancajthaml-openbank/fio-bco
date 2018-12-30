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
	"strconv"
	"strings"
)

// ReplyTimeout message
type ReplyTimeout struct{}

// TokenCreated message
type TokenCreated struct{}

// TokenDeleted message
type TokenDeleted struct{}

// Token represents metadata of token entity
type Token struct {
	Value        string `json:"value"`
	LastSyncedID int64  `json:"-"`
}

// Hydrate deserializes Token entity from persistent data
func (entity *Token) Hydrate(data []byte) {
	// FIXME more efficient read-split-inplace
	lines := strings.Split(string(data), "\n")

	if cast, err := strconv.ParseInt(lines[0], 10, 64); err == nil {
		entity.LastSyncedID = cast
	}
}
