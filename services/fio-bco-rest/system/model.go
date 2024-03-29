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

package system

// UnitStatus represents whitelist of properties of unit that we are willing to
// expose
type UnitStatus struct {
	Status          string `json:"status"`
	StatusChangedAt uint64 `json:"statusChangedAt"`
}

// CapacityStatus represents current capacity status
type CapacityStatus struct {
	Free      uint64 `json:"free"`
	Used      uint64 `json:"used"`
	IsHealthy bool   `json:"healthy"`
}

// Status represents system status snapshot
type Status struct {
	Storage CapacityStatus `json:"storage"`
	Memory  CapacityStatus `json:"memory"`
}
