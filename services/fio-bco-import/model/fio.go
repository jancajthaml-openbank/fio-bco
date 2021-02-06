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

// FioEnvelope represents fio gateway import statement envelope
type FioEnvelope struct {
	Info         FioAccountInfo
	Transactions []FioStatement
}

// FioAccountInfo represent chunk at accountStatement/info from transactions.json response
type FioAccountInfo struct {
	AccountID      string  `json:"accountId"`
	BankID         string  `json:"bankId"`
	Currency       string  `json:"currency"`
	IBAN           string  `json:"iban"`
	BIC            string  `json:"bic"`
	OpeningBalance float64 `json:"openingBalance"`
	ClosingBalance float64 `json:"closingBalance"`
	IDFrom         int     `json:"idFrom"`
	IDTo           int     `json:"idTo"`
	IDLastDownload int     `json:"idLastDownload"`
}

// FioStatement represent chunk of accountStatement/transactionList/transaction from transactions.json response
type FioStatement struct {
	TransferDate     *stringNode `json:"column0"`
	Amount           *floatNode  `json:"column1"`
	AccountTo        *stringNode `json:"column2"`
	AcountToBankCode *stringNode `json:"column3"`
	//TransferType  *stringNode `json:"column8"`  // FIXME e.g. "Příjem převodem uvnitř banky"
	Currency      *stringNode `json:"column14"`
	TransactionID *intNode    `json:"column17"`
	TransferID    *intNode    `json:"column22"`
	AccountToBIC  *stringNode `json:"column26"`
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
