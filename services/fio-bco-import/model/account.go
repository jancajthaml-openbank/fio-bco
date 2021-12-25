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

package model

// AccountVault entity
type AccountVault struct {
	Tenant string `json:"tenant"`
	Name   string `json:"name"`
}

// Account entity in vault format
type Account struct {
	Tenant         string `json:"-"`
	Name           string `json:"name"`
	Format         string `json:"format"`
	Currency       string `json:"currency"`
	IsBalanceCheck bool   `json:"isBalanceCheck"`
}
