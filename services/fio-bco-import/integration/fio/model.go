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

package fio

import (
	"encoding/json"
	"fmt"
)

type envelope struct {
	info       accountInfo
	statements []statement
}

type accountInfo struct {
	accountID string
	bankID    string
	currency  string
	iban      string
	bic       string
}

// UnmarshalJSON envelope
func (entity *envelope) UnmarshalJSON(data []byte) error {
	if entity == nil {
		return fmt.Errorf("cannot unmarshal to nil pointer")
	}

	all := struct {
		Statement struct {
			Info struct {
				AccountID string `json:"accountId"`
				BankID    string `json:"bankId"`
				Currency  string `json:"currency"`
				IBAN      string `json:"iban"`
				BIC       string `json:"bic"`
			} `json:"info"`
			TransactionList struct {
				Statements []statement `json:"transaction"`
			} `json:"transactionList"`
		} `json:"accountStatement"`
	}{}

	err := json.Unmarshal(data, &all)
	if err != nil {
		return err
	}
	if all.Statement.Info.AccountID == "" {
		return fmt.Errorf("missing attribute \"accountID\"")
	}
	if all.Statement.Info.BankID == "" {
		return fmt.Errorf("missing attribute \"bankId\"")
	}
	if all.Statement.Info.Currency == "" {
		return fmt.Errorf("missing attribute \"currency\"")
	}
	if all.Statement.Info.IBAN == "" {
		return fmt.Errorf("missing attribute \"iban\"")
	}
	if all.Statement.Info.BIC == "" {
		return fmt.Errorf("missing attribute \"bic\"")
	}

	entity.info.accountID = all.Statement.Info.AccountID
	entity.info.bankID = all.Statement.Info.BankID
	entity.info.currency = all.Statement.Info.Currency
	entity.info.iban = all.Statement.Info.IBAN
	entity.info.bic = all.Statement.Info.BIC
	entity.statements = all.Statement.TransactionList.Statements

	return nil
}

type statement struct {
	transferDate     *stringNode `json:"column0"`
	amount           *floatNode  `json:"column1"`
	accountTo        *stringNode `json:"column2"`
	acountToBankCode *stringNode `json:"column3"`
	accountToBIC     *stringNode `json:"column26"`
	//transferType     *stringNode `json:"column8"`  // FIXME e.g. "Příjem převodem uvnitř banky"
	currency         *stringNode `json:"column14"`
	transactionID    *intNode    `json:"column17"`
	transferID       *intNode    `json:"column22"`
}

type stringNode struct {
	value string `json:"value"`
	name  string `json:"name"`
	id    int    `json:"id"`
}

type intNode struct {
	value int64  `json:"value"`
	name  string `json:"name"`
	id    int    `json:"id"`
}

type floatNode struct {
	value float64 `json:"value"`
	name  string  `json:"name"`
	id    int     `json:"id"`
}
