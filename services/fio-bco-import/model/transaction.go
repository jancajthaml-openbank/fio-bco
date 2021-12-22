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

// Transaction entity in ledger-rest format
type Transaction struct {
	Tenant        string     `json:"-"`
	IDTransaction string     `json:"id"`
	Transfers     []Transfer `json:"transfers"`
}

// Transfer entity in ledger-rest format
type Transfer struct {
	ID int64 `json:"-"`
	IDTransfer string       `json:"id"`
	Credit     AccountVault `json:"credit"`
	Debit      AccountVault `json:"debit"`
	ValueDate  string       `json:"valueDate"`
	Amount     string       `json:"amount"`
	Currency   string       `json:"currency"`
}
