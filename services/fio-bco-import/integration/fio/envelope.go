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

//import (
//"encoding/json"
//"fmt"
//"github.com/jancajthaml-openbank/fio-bco-import/model"
//"math"
//"sort"
//"strconv"
//"time"
//)

// Envelope represents fio statements fascade
type Envelope struct {
	Info       Info
	Statements []Statement
}

// Info represents account information from which statements originate
type Info struct {
	AccountID string
	BankCode  string
	Currency  string
	IBAN      string
	BIC       string
}

/*
// GetTransactions returns transactions from fio statement
func (envelope *Envelope) GetTransactions(tenant string) []model.Transaction {
	transactions := make([]model.Transaction, 0)

	if envelope == nil {
		return transactions
	}

	sort.SliceStable(envelope.Statements, func(i, j int) bool {
		return envelope.Statements[i].TransactionID.Value < envelope.Statements[j].TransactionID.Value
	})

	previousIDTransaction := ""
	transfers := make([]model.Transfer, 0)

	now := time.Now()

	var credit string
	var debit string
	var currency string
	var valueDate time.Time

	for _, transfer := range envelope.Statements {
		if transfer.TransferID == nil || transfer.Amount == nil {
			continue
		}

		if transfer.Amount.Value > 0 {
			credit = envelope.Info.IBAN
			if transfer.AccountTo == nil {
				// INFO fee and taxes and maybe card payments
				debit = envelope.Info.BIC
			} else if transfer.AccountToBIC != nil {
				debit, _ = model.NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AccountToBIC.Value, "")
			} else if transfer.AcountToBankCode != nil {
				debit, _ = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", transfer.AcountToBankCode.Value)
			} else {
				debit, _ = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", envelope.Info.BankCode)
			}
		} else {
			if transfer.AccountTo == nil {
				// INFO fee and taxes and maybe card payments
				credit = envelope.Info.BIC
			} else if transfer.AccountToBIC != nil {
				credit, _ = model.NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AccountToBIC.Value, "")
			} else if transfer.AcountToBankCode != nil {
				credit, _ = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", transfer.AcountToBankCode.Value)
			} else {
				credit, _ = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", envelope.Info.BankCode)
			}
			debit = envelope.Info.IBAN
		}

		if transfer.TransferDate == nil {
			valueDate = now
		} else if date, err := time.Parse("2006-01-02-0700", transfer.TransferDate.Value); err == nil {
			valueDate = date.UTC()
		} else {
			valueDate = now
		}

		if transfer.Currency == nil {
			currency = envelope.Info.Currency
		} else {
			currency = transfer.Currency.Value
		}

		idTransaction := envelope.Info.IBAN + strconv.FormatInt(transfer.TransactionID.Value, 10)

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
			ID:         transfer.TransferID.Value,
			IDTransfer: strconv.FormatInt(transfer.TransferID.Value, 10),
			Credit: model.AccountVault{
				Tenant: tenant,
				Name:   credit,
			},
			Debit: model.AccountVault{
				Tenant: tenant,
				Name:   debit,
			},
			ValueDate: valueDate.Format("2006-01-02T15:04:05Z0700"),
			Amount:    strconv.FormatFloat(math.Abs(transfer.Amount.Value), 'f', -1, 64),
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

	set := make(map[string]model.Account)

	var normalizedAccount string
	var isIBAN bool
	var accountFormat string
	var currency string

	for _, transfer := range envelope.Statements {

		if transfer.AccountTo == nil {
			// INFO fee and taxes and maybe card payments
			normalizedAccount = envelope.Info.BIC
			isIBAN = false
		} else if transfer.AccountToBIC != nil {
			normalizedAccount, isIBAN = model.NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AccountToBIC.Value, "")
		} else if transfer.AcountToBankCode != nil {
			normalizedAccount, isIBAN = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", transfer.AcountToBankCode.Value)
		} else {
			normalizedAccount, isIBAN = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", envelope.Info.BankCode)
		}

		if transfer.AccountTo == nil {
			accountFormat = "FIO_TECHNICAL"
		} else if isIBAN {
			accountFormat = "IBAN"
			//fmt.Printf("Valid IBAN in statement: %s\n", normalizedAccount)
		} else {
			bytes, _ := json.Marshal(transfer)
			fmt.Printf("Strange account number in statement: %s\n", string(bytes))
			accountFormat = "FIO_UNKNOWN"
		}

		if transfer.Currency == nil {
			currency = envelope.Info.Currency
		} else {
			currency = transfer.Currency.Value
		}

		set[normalizedAccount] = model.Account{
			Tenant:         tenant,
			Name:           normalizedAccount,
			Format:         accountFormat,
			Currency:       currency,
			IsBalanceCheck: false,
		}

	}

	set[envelope.Info.IBAN] = model.Account{
		Tenant:         tenant,
		Name:           envelope.Info.IBAN,
		Format:         "IBAN",
		Currency:       envelope.Info.Currency,
		IsBalanceCheck: false,
	}

	for _, account := range set {
		accounts = append(accounts, account)
	}

	return accounts
}
*/
