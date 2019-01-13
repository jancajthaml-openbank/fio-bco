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
	"math"
	"sort"
	"strconv"
	"time"
)

type FioImportEnvelope struct {
	Statement FioAccountStatement `json:"accountStatement"`
}

type FioAccountStatement struct {
	Info            FioInfo            `json:"info"`
	TransactionList FioTransactionList `json:"transactionList"`
}

type FioInfo struct {
	AccountId      string  `json:"accountId"`
	BankId         string  `json:"bankId"`
	Currency       string  `json:"currency"`
	IBAN           string  `json:"iban"`
	BIC            string  `json:"bic"`
	OpeningBalance float64 `json:"openingBalance"`
	ClosingBalance float64 `json:"closingBalance"`
	IdFrom         int     `json:"idFrom"`
	IdTo           int     `json:"idTo"`
	IdLastDownload int     `json:"idLastDownload"`
}

type FioTransactionList struct {
	Transactions []FioTransaction `json:"transaction"`
}

type FioTransaction struct {
	Column0  *FioStringNode `json:"column0"`
	Column1  *FioFloatNode  `json:"column1"`
	Column2  *FioStringNode `json:"column2"`
	Column3  *FioStringNode `json:"column3"`
	Column4  *FioStringNode `json:"column4"`
	Column5  *FioStringNode `json:"column5"`
	Column6  *FioStringNode `json:"column6"`
	Column7  *FioStringNode `json:"column7"`
	Column8  *FioStringNode `json:"column8"`
	Column9  *FioStringNode `json:"column9"`
	Column10 *FioStringNode `json:"column10"`
	Column12 *FioStringNode `json:"column12"`
	Column14 *FioStringNode `json:"column14"`
	Column16 *FioStringNode `json:"column16"`
	Column17 *FioIntNode    `json:"column17"`
	Column18 *FioStringNode `json:"column18"`
	Column22 *FioIntNode    `json:"column22"`
	Column25 *FioStringNode `json:"column25"`
	Column26 *FioStringNode `json:"column26"`
}

type FioStringNode struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Id    int    `json:"id"`
}

type FioDateNode struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Id    int    `json:"id"`
}

type FioIntNode struct {
	Value int64  `json:"value"`
	Name  string `json:"name"`
	Id    int    `json:"id"`
}

type FioFloatNode struct {
	Value float64 `json:"value"`
	Name  string  `json:"name"`
	Id    int     `json:"id"`
}

func (envelope *FioImportEnvelope) GetTransactions() []Transaction {
	if envelope == nil {
		return nil
	}

	var set = make(map[int64][]Transfer)

	now := time.Now()

	var credit string
	var debit string
	var valueDate time.Time

	sort.SliceStable(envelope.Statement.TransactionList.Transactions, func(i, j int) bool {
		return envelope.Statement.TransactionList.Transactions[i].Column22.Value < envelope.Statement.TransactionList.Transactions[j].Column22.Value
	})

	for _, transfer := range envelope.Statement.TransactionList.Transactions {
		if transfer.Column22 == nil || transfer.Column1 == nil {
			continue
		}

		if transfer.Column1.Value > 0 {
			credit = envelope.Statement.Info.IBAN
			if transfer.Column2 == nil {
				debit = envelope.Statement.Info.BIC
			} else {
				if transfer.Column3 != nil {
					debit = NormalizeAccountNumber(transfer.Column2.Value, transfer.Column3.Value, envelope.Statement.Info.BankId)
				} else {
					debit = NormalizeAccountNumber(transfer.Column2.Value, "", envelope.Statement.Info.BankId)
				}
			}
		} else {
			if transfer.Column2 == nil {
				credit = envelope.Statement.Info.BIC
			} else {
				if transfer.Column3 != nil {
					credit = NormalizeAccountNumber(transfer.Column2.Value, transfer.Column3.Value, envelope.Statement.Info.BankId)
				} else {
					credit = NormalizeAccountNumber(transfer.Column2.Value, "", envelope.Statement.Info.BankId)
				}
			}
			debit = envelope.Statement.Info.IBAN
		}

		if transfer.Column0 == nil {
			valueDate = now
		} else if date, err := time.Parse("2006-01-02-0700", transfer.Column0.Value); err == nil {
			valueDate = date.UTC()
		} else {
			valueDate = now
		}

		set[transfer.Column17.Value] = append(set[transfer.Column17.Value], Transfer{
			IDTransfer: transfer.Column22.Value,
			Credit:     credit,
			Debit:      debit,
			ValueDate:  valueDate.Format("2006-01-02T15:04:05Z0700"), // FIXME not formatted with Z but instead with +100
			Amount:     math.Abs(transfer.Column1.Value),
			Currency:   envelope.Statement.Info.Currency, // FIXME not true in all cases
		})
	}

	result := make([]Transaction, 0)
	for transaction, transfers := range set {
		result = append(result, Transaction{
			IDTransaction: envelope.Statement.Info.IBAN + "/" + strconv.FormatInt(transaction, 10),
			Transfers:     transfers,
		})
	}

	return result

}

func (envelope *FioImportEnvelope) GetAccounts() []Account {
	if envelope == nil {
		return nil
	}

	var set = make(map[string]FioTransaction)

	for _, transfer := range envelope.Statement.TransactionList.Transactions {
		if transfer.Column2 == nil {
			// INFO fee and taxes and maybe card payments
			set[envelope.Statement.Info.BIC] = transfer
		} else {
			set[transfer.Column2.Value] = transfer
		}
	}

	var normalizedAccount string
	var deduplicated = make(map[string]Account)

	for account, transfer := range set {
		if transfer.Column3 != nil {
			normalizedAccount = NormalizeAccountNumber(account, transfer.Column3.Value, envelope.Statement.Info.BankId)
		} else {
			normalizedAccount = NormalizeAccountNumber(account, "", envelope.Statement.Info.BankId)
		}

		deduplicated[normalizedAccount] = Account{
			Name:           normalizedAccount,
			Currency:       envelope.Statement.Info.Currency, // FIXME not true in all cases
			IsBalanceCheck: false,
		}
	}

	deduplicated[envelope.Statement.Info.IBAN] = Account{
		Name:           envelope.Statement.Info.IBAN,
		Currency:       envelope.Statement.Info.Currency,
		IsBalanceCheck: false,
	}

	result := make([]Account, 0)
	for _, item := range deduplicated {
		result = append(result, item)
	}

	return result
}

// http://cavaliercoder.com/blog/optimized-abs-for-int64-in-go.html
