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

package integration

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/metrics"
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"
	"github.com/jancajthaml-openbank/fio-bco-import/support/http"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// Workflow represents import integration workflow
type Workflow struct {
	Token            *model.Token
	Tenant           string
	FioClient        *http.FioClient
	VaultClient      *http.VaultClient
	LedgerClient     *http.LedgerClient
	EncryptedStorage localfs.Storage
	PlaintextStorage localfs.Storage
	Metrics          metrics.Metrics
}

// NewWorkflow returns fascade for integration workflow
func NewWorkflow(
	token *model.Token,
	tenant string,
	fioGateway string,
	vaultGateway string,
	ledgerGateway string,
	encryptedStorage localfs.Storage,
	plaintextStorage localfs.Storage,
	metrics metrics.Metrics,
) Workflow {
	return Workflow{
		Token:            token,
		Tenant:           tenant,
		FioClient:        http.NewFioClient(fioGateway),
		VaultClient:      http.NewVaultClient(vaultGateway),
		LedgerClient:     http.NewLedgerClient(ledgerGateway),
		EncryptedStorage: encryptedStorage,
		PlaintextStorage: plaintextStorage,
		Metrics:          metrics,
	}
}

func createAccountsFromStatements(
	tenant string,
	vaultClient *http.VaultClient,
	envelope *model.FioEnvelope,
) error {
	if envelope == nil {
		return fmt.Errorf("nil statements")
	}

	accounts := make(map[string]model.Account)

	var set = make(map[string]model.FioStatement)

	for _, transfer := range envelope.Transactions {
		if transfer.AccountTo == nil {
			// INFO fee and taxes and maybe card payments
			set[envelope.Info.BIC] = transfer
		} else {
			set[transfer.AccountTo.Value] = transfer
		}
	}

	var normalizedAccount string
	var accountFormat string
	var currency string

	for account, transfer := range set {
		if transfer.AcountToBankCode != nil {
			normalizedAccount = model.NormalizeAccountNumber(account, transfer.AcountToBankCode.Value, envelope.Info.BankID)
		} else {
			normalizedAccount = model.NormalizeAccountNumber(account, "", envelope.Info.BankID)
		}

		if normalizedAccount != account {
			accountFormat = "IBAN"
		} else {
			accountFormat = "FIO_UNKNOWN"
		}

		if transfer.Currency == nil {
			currency = envelope.Info.Currency
		} else {
			currency = transfer.Currency.Value
		}

		accounts[normalizedAccount] = model.Account{
			Tenant:         tenant,
			Name:           normalizedAccount,
			Format:         accountFormat,
			Currency:       currency,
			IsBalanceCheck: false,
		}

	}

	accounts[envelope.Info.IBAN] = model.Account{
		Tenant:         tenant,
		Name:           envelope.Info.IBAN,
		Format:         "IBAN",
		Currency:       envelope.Info.Currency,
		IsBalanceCheck: false,
	}

	for _, account := range accounts {
		log.Info().Msgf("Creating account %s", account.Name)

		err := vaultClient.CreateAccount(account)
		if err != nil {
			return fmt.Errorf("unable to create account %s", account.Name)
		}
	}

	return nil
}

func createTransactionsFromStatements(
	tenant string,
	ledgerClient *http.LedgerClient,
	encryptedStorage localfs.Storage,
	metrics metrics.Metrics,
	token *model.Token,
	envelope *model.FioEnvelope,
) error {
	if envelope == nil {
		return fmt.Errorf("nil statements")
	}

	sort.SliceStable(envelope.Transactions, func(i, j int) bool {
		return envelope.Transactions[i].TransactionID.Value < envelope.Transactions[j].TransactionID.Value
	})

	previousIDTransaction := ""
	transfers := make([]model.Transfer, 0)

	now := time.Now()

	var credit string
	var debit string
	var currency string
	var valueDate time.Time

	for _, transfer := range envelope.Transactions {
		if transfer.TransferID == nil || transfer.Amount == nil {
			continue
		}

		if transfer.Amount.Value > 0 {
			credit = envelope.Info.IBAN
			if transfer.AccountTo == nil {
				debit = envelope.Info.BIC
			} else {
				if transfer.AcountToBankCode != nil {
					debit = model.NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AcountToBankCode.Value, envelope.Info.BankID)
				} else if transfer.AccountToBIC != nil {
					debit = model.NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AccountToBIC.Value, envelope.Info.BankID)
				} else {
					debit = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", envelope.Info.BankID)
				}
			}
		} else {
			if transfer.AccountTo == nil {
				credit = envelope.Info.BIC
			} else {
				if transfer.AcountToBankCode != nil {
					credit = model.NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AcountToBankCode.Value, envelope.Info.BankID)
				} else if transfer.AccountToBIC != nil {
					credit = model.NormalizeAccountNumber(transfer.AccountTo.Value, transfer.AccountToBIC.Value, envelope.Info.BankID)
				} else {
					credit = model.NormalizeAccountNumber(transfer.AccountTo.Value, "", envelope.Info.BankID)
				}
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
			log.Info().Msgf("Creating transaction %s", previousIDTransaction)
			err := ledgerClient.CreateTransaction(model.Transaction{
				Tenant:        tenant,
				IDTransaction: previousIDTransaction,
				Transfers:     transfers,
			})
			if err != nil {
				return fmt.Errorf("unable to create transaction %s/%s", tenant, previousIDTransaction)
			}
			metrics.TransactionImported(len(transfers))
			for _, transfer := range transfers {
				if token.LastSyncedID > transfer.IDTransfer {
					continue
				}
				token.LastSyncedID = transfer.IDTransfer
				if !persistence.UpdateToken(encryptedStorage, token) {
					log.Warn().Msgf("unable to update token %s", token.ID)
				}
			}
			previousIDTransaction = idTransaction
			transfers = make([]model.Transfer, 0)
		}

		transfers = append(transfers, model.Transfer{
			IDTransfer: transfer.TransferID.Value,
			Credit: model.AccountPair{
				Tenant: tenant,
				Name:   credit,
			},
			Debit: model.AccountPair{
				Tenant: tenant,
				Name:   debit,
			},
			ValueDate: valueDate.Format("2006-01-02T15:04:05Z0700"),
			Amount:    strconv.FormatFloat(math.Abs(transfer.Amount.Value), 'f', -1, 64),
			Currency:  currency,
		})
	}

	if len(transfers) == 0 {
		return nil
	}

	log.Info().Msgf("Creating transaction %s", previousIDTransaction)
	err := ledgerClient.CreateTransaction(model.Transaction{
		Tenant:        tenant,
		IDTransaction: previousIDTransaction,
		Transfers:     transfers,
	})
	if err != nil {
		return fmt.Errorf("unable to create transaction %s/%s", tenant, previousIDTransaction)
	}
	metrics.TransactionImported(len(transfers))
	for _, transfer := range transfers {
		if token.LastSyncedID > transfer.IDTransfer {
			continue
		}
		token.LastSyncedID = transfer.IDTransfer
		if !persistence.UpdateToken(encryptedStorage, token) {
			log.Warn().Msgf("unable to update token %s", token.ID)
		}
	}

	return nil
}

func synchronizeNewStatements(
	encryptedStorage localfs.Storage,
	plaintextStorage localfs.Storage,
	token *model.Token,
	tenant string,
	fioClient *http.FioClient,
	vaultClient *http.VaultClient,
	ledgerClient *http.LedgerClient,
	metrics metrics.Metrics,
) {
	if token == nil {
		return
	}

	envelope, err := fioClient.GetStatementsEnvelope(*token)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to get envelope")
		return
	}
	if len(envelope.Transactions) == 0 {
		return
	}

	log.Debug().Msgf("token %s importing accounts", token.ID)
	err = createAccountsFromStatements(tenant, vaultClient, envelope)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to create accounts from envelope")
		return
	}

	log.Debug().Msgf("token %s importing transactions", token.ID)
	err = createTransactionsFromStatements(tenant, ledgerClient, encryptedStorage, metrics, token, envelope)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to create transactions from envelope")
		return
	}
}

// SynchronizeStatements downloads new statements from fio gateway and creates accounts and transactions in core
func (workflow Workflow) SynchronizeStatements() {
	synchronizeNewStatements(
		workflow.EncryptedStorage,
		workflow.PlaintextStorage,
		workflow.Token,
		workflow.Tenant,
		workflow.FioClient,
		workflow.VaultClient,
		workflow.LedgerClient,
		workflow.Metrics,
	)
}
