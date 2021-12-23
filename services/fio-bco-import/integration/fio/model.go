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

// UnmarshalJSON Envelope
func (entity *Envelope) UnmarshalJSON(data []byte) error {
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
				Statements []Statement `json:"transaction"`
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

	entity.AccountID = all.Statement.Info.AccountID
	entity.BankID = all.Statement.Info.BankID
	entity.Currency = all.Statement.Info.Currency
	entity.IBAN = all.Statement.Info.IBAN
	entity.BIC = all.Statement.Info.BIC
	entity.Statements = all.Statement.TransactionList.Statements

	return nil
}

type Statement struct {
	TransferDate     *stringNode `json:"column0"`
	Amount           *floatNode  `json:"column1"`
	AccountTo        *stringNode `json:"column2"`
	AcountToBankCode *stringNode `json:"column3"`
	AccountToBIC     *stringNode `json:"column26"`
	//TransferType     *stringNode `json:"column8"`  // FIXME e.g. "Příjem převodem uvnitř banky"
	Currency      *stringNode `json:"column14"`
	TransactionID *intNode    `json:"column17"`
	TransferID    *intNode    `json:"column22"`
}

type stringNode struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}

type intNode struct {
	Value int64  `json:"value"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}

type floatNode struct {
	Value float64 `json:"value"`
	Name  string  `json:"name"`
	ID    int     `json:"id"`
}
