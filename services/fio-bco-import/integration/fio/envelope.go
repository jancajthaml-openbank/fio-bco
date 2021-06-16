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
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"math"
	"sort"
	"strconv"
	"time"
)

// Envelope represents fio statements fascade
type Envelope struct {
	data envelope
}

// NewEnvelope returns fio statement accessor
func NewEnvelope(data envelope) *Envelope {
	return &Envelope{
		data: data,
	}
}

// GetTransactions returns transactions from fio statement
func (envelope *Envelope) GetTransactions(tenant string) []model.Transaction {
	transactions := make([]model.Transaction, 0)

	if envelope == nil {
		return transactions
	}

	sort.SliceStable(envelope.data.statements, func(i, j int) bool {
		return envelope.data.statements[i].transactionID.value < envelope.data.statements[j].transactionID.value
	})

	previousIDTransaction := ""
	transfers := make([]model.Transfer, 0)

	now := time.Now()

	var credit string
	var debit string
	var currency string
	var valueDate time.Time

	for _, transfer := range envelope.data.statements {
		if transfer.transferID == nil || transfer.amount == nil {
			continue
		}

		if transfer.amount.value > 0 {
			credit = envelope.data.info.iban
			if transfer.accountTo == nil {
				debit = envelope.data.info.bic
			} else {
				if transfer.acountToBankCode != nil {
					debit = model.NormalizeAccountNumber(transfer.accountTo.value, transfer.acountToBankCode.value, envelope.data.info.bankID)
				} else if transfer.accountToBIC != nil {
					debit = model.NormalizeAccountNumber(transfer.accountTo.value, transfer.accountToBIC.value, envelope.data.info.bankID)
				} else {
					debit = model.NormalizeAccountNumber(transfer.accountTo.value, "", envelope.data.info.bankID)
				}
			}
		} else {
			if transfer.accountTo == nil {
				credit = envelope.data.info.bic
			} else {
				if transfer.acountToBankCode != nil {
					credit = model.NormalizeAccountNumber(transfer.accountTo.value, transfer.acountToBankCode.value, envelope.data.info.bankID)
				} else if transfer.accountToBIC != nil {
					credit = model.NormalizeAccountNumber(transfer.accountTo.value, transfer.accountToBIC.value, envelope.data.info.bankID)
				} else {
					credit = model.NormalizeAccountNumber(transfer.accountTo.value, "", envelope.data.info.bankID)
				}
			}
			debit = envelope.data.info.iban
		}

		if transfer.transferDate == nil {
			valueDate = now
		} else if date, err := time.Parse("2006-01-02-0700", transfer.transferDate.value); err == nil {
			valueDate = date.UTC()
		} else {
			valueDate = now
		}

		if transfer.currency == nil {
			currency = envelope.data.info.currency
		} else {
			currency = transfer.currency.value
		}

		idTransaction := envelope.data.info.iban + strconv.FormatInt(transfer.transactionID.value, 10)

		if previousIDTransaction == "" {
			previousIDTransaction = idTransaction
		} else if previousIDTransaction != idTransaction {
			transactions = append(transactions, model.Transaction{
				Tenant:        tenant,
				IDTransaction: previousIDTransaction,
				Transfers:     transfers,
			})
			previousIDTransaction = idTransaction
			transfers = make([]model.Transfer, 0)
		}

		transfers = append(transfers, model.Transfer{
			IDTransfer: transfer.transferID.value,
			Credit: model.AccountPair{
				Tenant: tenant,
				Name:   credit,
			},
			Debit: model.AccountPair{
				Tenant: tenant,
				Name:   debit,
			},
			ValueDate: valueDate.Format("2006-01-02T15:04:05Z0700"),
			Amount:    strconv.FormatFloat(math.Abs(transfer.amount.value), 'f', -1, 64),
			Currency:  currency,
		})
	}

	if len(transfers) != 0 {
		transactions = append(transactions, model.Transaction{
			Tenant:        tenant,
			IDTransaction: previousIDTransaction,
			Transfers:     transfers,
		})
	}

	return transactions
}

// GetAccounts returns accounts from fio statement
func (envelope *Envelope) GetAccounts(tenant string) []model.Account {

	accounts := make([]model.Account, 0)

	if envelope == nil {
		return accounts
	}

	statements := make(map[string]statement)

	for _, transfer := range envelope.data.statements {
		if transfer.accountTo == nil {
			// INFO fee and taxes and maybe card payments
			statements[envelope.data.info.bic] = transfer
		} else {
			statements[transfer.accountTo.value] = transfer
		}
	}

	set := make(map[string]model.Account)

	var normalizedAccount string
	var accountFormat string
	var currency string

	for account, transfer := range statements {
		if transfer.acountToBankCode != nil {
			normalizedAccount = model.NormalizeAccountNumber(account, transfer.acountToBankCode.value, envelope.data.info.bankID)
		} else {
			normalizedAccount = model.NormalizeAccountNumber(account, "", envelope.data.info.bankID)
		}

		if normalizedAccount != account {
			accountFormat = "IBAN"
		} else {
			accountFormat = "FIO_UNKNOWN"
		}

		if transfer.currency == nil {
			currency = envelope.data.info.currency
		} else {
			currency = transfer.currency.value
		}

		set[normalizedAccount] = model.Account{
			Tenant:         tenant,
			Name:           normalizedAccount,
			Format:         accountFormat,
			Currency:       currency,
			IsBalanceCheck: false,
		}

	}

	set[envelope.data.info.iban] = model.Account{
		Tenant:         tenant,
		Name:           envelope.data.info.iban,
		Format:         "IBAN",
		Currency:       envelope.data.info.currency,
		IsBalanceCheck: false,
	}

	for _, account := range set {
		accounts = append(accounts, account)
	}

	return accounts
}
