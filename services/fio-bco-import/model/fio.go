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
	"math"
	"strconv"
	"time"
)

// ImportEnvelope represents fio gateway import statement entity
type ImportEnvelope struct {
	Statement accountStatement `json:"accountStatement"`
}

type accountStatement struct {
	Info            accountInfo     `json:"info"`
	TransactionList transactionList `json:"transactionList"`
}

type accountInfo struct {
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

type transactionList struct {
	Transactions []fioTransaction `json:"transaction"`
}

type fioTransaction struct {
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

type dateNode struct {
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

// GetTransactions return generator of fio transactions over given envelope
func (envelope *ImportEnvelope) GetTransactions(tenant string) <-chan Transaction {
	chnl := make(chan Transaction)
	if envelope == nil {
		close(chnl)
		return chnl
	}

	var previousIDTransaction = ""
	var buffer = make([]Transfer, 0)

	go func() {
		defer close(chnl)

		now := time.Now()

		var credit string
		var debit string
		var currency string
		var valueDate time.Time

		for _, transfer := range envelope.Statement.TransactionList.Transactions {
			if transfer.TransferID == nil || transfer.Amount == nil {
				continue
			}

			if transfer.Amount.Value > 0 {
				credit = envelope.Statement.Info.IBAN
				if transfer.AccountTo == nil {
					debit = envelope.Statement.Info.BIC
				} else {
					if transfer.AcountToBankCode != nil {
						debit = NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AcountToBankCode.Value, envelope.Statement.Info.BankID)
					} else if transfer.AccountToBIC != nil {
						debit = NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AccountToBIC.Value, envelope.Statement.Info.BankID)
					} else {
						debit = NormalizeAccountNumber(transfer.AccountTo.Value, "", envelope.Statement.Info.BankID)
					}
				}
			} else {
				if transfer.AccountTo == nil {
					credit = envelope.Statement.Info.BIC
				} else {
					if transfer.AcountToBankCode != nil {
						credit = NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AcountToBankCode.Value, envelope.Statement.Info.BankID)
					} else if transfer.AccountToBIC != nil {
						credit = NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AccountToBIC.Value, envelope.Statement.Info.BankID)
					} else {
						credit = NormalizeAccountNumber(transfer.AccountTo.Value, "", envelope.Statement.Info.BankID)
					}
				}
				debit = envelope.Statement.Info.IBAN
			}

			if transfer.TransferDate == nil {
				valueDate = now
			} else if date, err := time.Parse("2006-01-02-0700", transfer.TransferDate.Value); err == nil {
				valueDate = date.UTC()
			} else {
				valueDate = now
			}

			if transfer.Currency == nil {
				currency = envelope.Statement.Info.Currency
			} else {
				currency = transfer.Currency.Value
			}

			idTransaction := envelope.Statement.Info.IBAN + strconv.FormatInt(transfer.TransactionID.Value, 10)

			if previousIDTransaction == "" {
				previousIDTransaction = idTransaction
			} else if previousIDTransaction != idTransaction {
				transfers := make([]Transfer, len(buffer))
				copy(transfers, buffer)
				buffer = make([]Transfer, 0)
				chnl <- Transaction{
					Tenant:        tenant,
					IDTransaction: previousIDTransaction,
					Transfers:     transfers,
				}
				previousIDTransaction = idTransaction
			}

			buffer = append(buffer, Transfer{
				IDTransfer: transfer.TransferID.Value,
				Credit: AccountPair{
					Tenant: tenant,
					Name:   credit,
				},
				Debit: AccountPair{
					Tenant: tenant,
					Name:   debit,
				},
				ValueDate: valueDate.Format("2006-01-02T15:04:05Z0700"),
				Amount:    strconv.FormatFloat(math.Abs(transfer.Amount.Value), 'f', -1, 64),
				Currency:  currency,
			})
		}

		if len(buffer) == 0 {
			return
		}

		transfers := make([]Transfer, len(buffer))
		copy(transfers, buffer)
		buffer = make([]Transfer, 0)
		chnl <- Transaction{
			Tenant:        tenant,
			IDTransaction: previousIDTransaction,
			Transfers:     transfers,
		}
	}()

	return chnl
}

// GetAccounts return generator of fio accounts over given envelope
func (envelope *ImportEnvelope) GetAccounts(tenant string) <-chan Account {
	chnl := make(chan Account)
	if envelope == nil {
		close(chnl)
		return chnl
	}

	var visited = make(map[string]interface{})

	go func() {
		defer close(chnl)

		var set = make(map[string]fioTransaction)

		for _, transfer := range envelope.Statement.TransactionList.Transactions {
			if transfer.AccountTo == nil {
				// INFO fee and taxes and maybe card payments
				set[envelope.Statement.Info.BIC] = transfer
			} else {
				set[transfer.AccountTo.Value] = transfer
			}
		}

		var normalizedAccount string
		var accountFormat string
		var currency string

		for account, transfer := range set {
			if transfer.AcountToBankCode != nil {
				normalizedAccount = NormalizeAccountNumber(account, transfer.AcountToBankCode.Value, envelope.Statement.Info.BankID)
			} else {
				normalizedAccount = NormalizeAccountNumber(account, "", envelope.Statement.Info.BankID)
			}

			if normalizedAccount != account {
				accountFormat = "IBAN"
			} else {
				accountFormat = "FIO_UNKNOWN"
			}

			if transfer.Currency == nil {
				currency = envelope.Statement.Info.Currency
			} else {
				currency = transfer.Currency.Value
			}

			if _, ok := visited[normalizedAccount]; !ok {
				chnl <- Account{
					Tenant:         tenant,
					Name:           normalizedAccount,
					Format:         accountFormat,
					Currency:       currency,
					IsBalanceCheck: false,
				}
				visited[normalizedAccount] = nil
			}
		}

		if _, ok := visited[envelope.Statement.Info.IBAN]; !ok {
			chnl <- Account{
				Tenant:         tenant,
				Name:           envelope.Statement.Info.IBAN,
				Format:         "IBAN",
				Currency:       envelope.Statement.Info.Currency,
				IsBalanceCheck: false,
			}
			visited[envelope.Statement.Info.IBAN] = nil
		}
	}()

	return chnl
}
