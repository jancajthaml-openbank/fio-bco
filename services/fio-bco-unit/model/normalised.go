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
	"github.com/jancajthaml-openbank/fio-bco-unit/iban"
)

type Transaction struct {
	IDTransaction int64      `json:"id,string"`
	Transfers     []Transfer `json:"transfers"`
}

type Transfer struct {
	IDTransfer int64   `json:"id,string"`
	Credit     string  `json:"credit"`
	Debit      string  `json:"debit"`
	ValueDate  string  `json:"valueDate"`
	Amount     float64 `json:"amount,string"`
	Currency   string  `json:"currency"`
}

type Account struct {
	Name           string `json:"accountNumber"`
	Currency       string `json:"currency"`
	IsBalanceCheck bool   `json:"isBalanceCheck"`
}

func NormalizeAccountNumber(number string, bankCode string, nostroBankCode string) string {
	var calculatedIBAN string

	if bankCode == "" {
		calculatedIBAN = iban.Calculate(number, nostroBankCode)
	} else {
		calculatedIBAN = iban.Calculate(number, bankCode)
	}

	if calculatedIBAN == "" {
		return number
	}

	return calculatedIBAN
}
