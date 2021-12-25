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
	"encoding/json"
	"fmt"
	"time"
	"math"
	"strconv"

	"github.com/jancajthaml-openbank/fio-bco-import/integration/fio"
	"github.com/jancajthaml-openbank/fio-bco-import/integration/ledger"
	"github.com/jancajthaml-openbank/fio-bco-import/integration/vault"
	"github.com/jancajthaml-openbank/fio-bco-import/metrics"
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/persistence"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// Workflow represents import integration workflow
type Workflow struct {
	Token            *model.Token
	Tenant           string
	FioClient        *fio.Client
	VaultClient      *vault.Client
	LedgerClient     *ledger.Client
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
		FioClient:        fio.NewClient(fioGateway),
		VaultClient:      vault.NewClient(vaultGateway),
		LedgerClient:     ledger.NewClient(ledgerGateway),
		EncryptedStorage: encryptedStorage,
		PlaintextStorage: plaintextStorage,
		Metrics:          metrics,
	}
}

// DownloadStatements download new statements from fio gateway
func (workflow Workflow) DownloadStatements() {
	if workflow.Token == nil {
		return
	}

	log.Debug().Msgf("token %s downloading new statements", workflow.Token.ID)

	envelope, err := workflow.FioClient.GetStatementsEnvelope(*workflow.Token)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to get envelope")
		return
	}

	exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/nostro")
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to check if info for token %sexists with error %s", workflow.Token.ID, err)
		return
	}

	if !exists {
		log.Info().Msgf("Should create info about synchronized account before downloading statements")
		data, err := json.Marshal(envelope.Info)
		if err != nil {
			log.Warn().Msgf("Unable to marshal info of %s", workflow.Token.ID)
			return
		}
		err = workflow.PlaintextStorage.WriteFileExclusive("token/"+workflow.Token.ID+"/nostro", data)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to persist envelope nostro info of %s with error %s", workflow.Token.ID, err)
			return
		}
	}

	log.Debug().Msgf("token %s downloaded %d new statements", workflow.Token.ID, len(envelope.Statements))

	for _, transfer := range envelope.Statements {
		if transfer.TransferID == nil {
			continue
		}
		exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/statements/" + strconv.FormatInt(transfer.TransferID.Value, 10))
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to check if transaction %d exists for token %s with error %s", transfer.TransferID.Value, workflow.Token.ID, err)
			return
		}
		if exists {
			continue
		}
		data, err := json.Marshal(transfer)
		if err != nil {
			log.Warn().Msgf("Unable to marshal statement details of %s/%d", workflow.Token.ID, transfer.TransferID.Value)
			continue
		}
		err = workflow.PlaintextStorage.WriteFileExclusive("token/"+workflow.Token.ID+"/statements/"+strconv.FormatInt(transfer.TransferID.Value, 10)+"/data", data)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to persist statement details of %s/%s/%d with error %s", workflow.Token.ID, transfer.TransferID.Value, err)
			continue
		}
		if workflow.Token.LastSyncedID >= transfer.TransferID.Value {
			continue
		}
		workflow.Token.LastSyncedID = transfer.TransferID.Value
		if !persistence.UpdateToken(workflow.EncryptedStorage, workflow.Token) {
			log.Warn().Msgf("unable to update token %s", workflow.Token.ID)
		}
	}
}

func (workflow Workflow) CreateAccounts() {
	log.Debug().Msgf("token %s creating accounts from statements", workflow.Token.ID)

	data, err := workflow.PlaintextStorage.ReadFileFully("token/" + workflow.Token.ID + "/nostro")
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to load %s nostro info", workflow.Token.ID)
		return
	}

	info := new(fio.Info)
	if json.Unmarshal(data, info) != nil {
		log.Warn().Msgf("Unable to unmarshal info %s", workflow.Token.ID)
		return
	}

	set := make(map[string]model.Account)
	idsNeedingConfirmation := make([]string, 0)
	creatingNostro := false

	exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/ack_nostro")
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to check if %s nostro account ack exists exists", workflow.Token.ID)
		return
	}

	if !exists {
		set[info.IBAN] = model.Account{
			Tenant:         workflow.Tenant,
			Name:           info.IBAN,
			Format:         "IBAN",
			Currency:       info.Currency,
			IsBalanceCheck: false,
		}
		creatingNostro = true
	}

	ids, err := workflow.PlaintextStorage.ListDirectory("token/"+workflow.Token.ID+"/statements", true)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to obtain transaction ids from storage for token %s", workflow.Token.ID)
		return
	}

	for _, id := range ids {
		exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/statements/" + id + "/ack_account")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to check if statement %s/%s accounts exists", workflow.Token.ID, id)
			continue
		}
		if exists {
			continue
		}

		data, err := workflow.PlaintextStorage.ReadFileFully("token/" + workflow.Token.ID + "/statements/" + id + "/data")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to load statement %s/%s", workflow.Token.ID, id)
			continue
		}

		statement := new(fio.Statement)
		if json.Unmarshal(data, statement) != nil {
			log.Warn().Msgf("Unable to unmarshal statement %s/%s", workflow.Token.ID, id)
			continue
		}

		idsNeedingConfirmation = append(idsNeedingConfirmation, id)

		var normalizedAccount string
		var isIBAN bool
		var accountFormat string
		var currency string

		if statement.AccountTo == nil {
			// INFO fee and taxes and maybe card payments
			normalizedAccount = info.BIC + "_" + info.Currency
			isIBAN = false
		} else if statement.AccountToBIC != nil {
			normalizedAccount, isIBAN = model.NormalizeAccountNumber(statement.AccountTo.Value, statement.AccountToBIC.Value, "")
		} else if statement.AcountToBankCode != nil {
			normalizedAccount, isIBAN = model.NormalizeAccountNumber(statement.AccountTo.Value, "", statement.AcountToBankCode.Value)
		} else {
			normalizedAccount, isIBAN = model.NormalizeAccountNumber(statement.AccountTo.Value, "", info.BankCode)
		}

		if statement.AccountTo == nil {
			accountFormat = "FIO_TECHNICAL"
		} else if isIBAN {
			accountFormat = "IBAN"
		} else {
			fmt.Printf("Strange account number in statement: %s\n", string(data))
			accountFormat = "FIO_UNKNOWN"
		}

		if statement.Specification != nil {
			currency = statement.Specification.Value[len(statement.Specification.Value)-3 : len(statement.Specification.Value)]
		} else if statement.Currency != nil {
			currency = statement.Currency.Value
		} else {
			currency = info.Currency
		}

		if statement.Specification != nil {
			set[info.BIC+"_"+info.Currency] = model.Account{
				Tenant:         workflow.Tenant,
				Name:           info.BIC+"_"+info.Currency,
				Format:         "FIO_TECHNICAL",
				Currency:       info.Currency,
				IsBalanceCheck: false,
			}
			set[info.BIC+"_"+currency] = model.Account{
				Tenant:         workflow.Tenant,
				Name:           info.BIC+"_"+currency,
				Format:         "FIO_TECHNICAL",
				Currency:       currency,
				IsBalanceCheck: false,
			}
		}

		set[normalizedAccount] = model.Account{
			Tenant:         workflow.Tenant,
			Name:           normalizedAccount,
			Format:         accountFormat,
			Currency:       currency,
			IsBalanceCheck: false,
		}
	}

	for _, account := range set {
		log.Info().Msgf("Creating account %s", account.Name)

		err := workflow.VaultClient.CreateAccount(account)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to create account %s", account.Name)
			return
		}
	}

	if creatingNostro {
		err = workflow.PlaintextStorage.TouchFile("token/" + workflow.Token.ID + "/ack_nostro")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to mark nostro account as createdfor %s", workflow.Token.ID)
		}
	}

	for _, id := range idsNeedingConfirmation {
		err = workflow.PlaintextStorage.TouchFile("token/" + workflow.Token.ID + "/statements/" + id + "/ack_account")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to mark account discovery for %s/%s", workflow.Token.ID, id)
		}
	}
}


func (workflow Workflow) CreateTransactions() {
	log.Debug().Msgf("token %s creating transactions from statements", workflow.Token.ID)

	data, err := workflow.PlaintextStorage.ReadFileFully("token/" + workflow.Token.ID + "/nostro")
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to load %s nostro info", workflow.Token.ID)
		return
	}

	info := new(fio.Info)
	if json.Unmarshal(data, info) != nil {
		log.Warn().Msgf("Unable to unmarshal info %s", workflow.Token.ID)
		return
	}

	ids, err := workflow.PlaintextStorage.ListDirectory("token/"+workflow.Token.ID+"/statements", true)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to obtain transaction ids from storage for token %s", workflow.Token.ID)
		return
	}

	previousIDTransaction := ""
	idsNeedingConfirmation := make([]string, 0)
	transfers := make([]model.Transfer, 0)

	now := time.Now()

	var credit string
	var debit string
	var creditCurrency string
	var creditAmount string
	var debitCurrency string
	var debitAmount string
	var valueDate time.Time

	for _, id := range ids {
		exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/statements/" + id + "/ack_transfer")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to check if statement %s/%s transfer exists", workflow.Token.ID, id)
			continue
		}
		if exists {
			continue
		}

		data, err := workflow.PlaintextStorage.ReadFileFully("token/" + workflow.Token.ID + "/statements/" + id + "/data")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to load statement %s/%s", workflow.Token.ID, id)
			continue
		}

		statement := new(fio.Statement)
		if json.Unmarshal(data, statement) != nil {
			log.Warn().Msgf("Unable to unmarshal statement %s/%s", workflow.Token.ID, id)
			continue
		}

		if statement.TransferID == nil || statement.Amount == nil {
			continue
		}

		if statement.Amount.Value > 0 {
			credit = info.IBAN
			if statement.AccountTo == nil {
				// INFO fee and taxes and maybe card payments
				debit = info.BIC + "_" + info.Currency
			} else if statement.AccountToBIC != nil {
				debit, _ = model.NormalizeAccountNumber(statement.AccountTo.Value, statement.AccountToBIC.Value, "")
			} else if statement.AcountToBankCode != nil {
				debit, _ = model.NormalizeAccountNumber(statement.AccountTo.Value, "", statement.AcountToBankCode.Value)
			} else {
				debit, _ = model.NormalizeAccountNumber(statement.AccountTo.Value, "", info.BankCode)
			}
		} else {
			if statement.AccountTo == nil {
				// INFO fee and taxes and maybe card payments
				credit = info.BIC + "_" + info.Currency
			} else if statement.AccountToBIC != nil {
				credit, _ = model.NormalizeAccountNumber(statement.AccountTo.Value, statement.AccountToBIC.Value, "")
			} else if statement.AcountToBankCode != nil {
				credit, _ = model.NormalizeAccountNumber(statement.AccountTo.Value, "", statement.AcountToBankCode.Value)
			} else {
				credit, _ = model.NormalizeAccountNumber(statement.AccountTo.Value, "", info.BankCode)
			}
			debit = info.IBAN
		}

		if statement.TransferDate == nil {
			valueDate = now
		} else if date, err := time.Parse("2006-01-02-0700", statement.TransferDate.Value); err == nil {
			valueDate = date.UTC()
		} else {
			valueDate = now
		}

		if statement.Specification != nil && statement.Currency != nil {
			creditCurrency = statement.Specification.Value[len(statement.Specification.Value)-3 : len(statement.Specification.Value)]
			creditAmount = statement.Specification.Value[0: len(statement.Specification.Value) - 4]
			debitCurrency = statement.Currency.Value
			debitAmount = strconv.FormatFloat(math.Abs(statement.Amount.Value), 'f', -1, 64)
		} else if statement.Specification != nil {
			creditCurrency = statement.Specification.Value[len(statement.Specification.Value)-3 : len(statement.Specification.Value)]
			creditAmount = statement.Specification.Value[0: len(statement.Specification.Value) - 4]
			debitCurrency = info.Currency
			debitAmount = strconv.FormatFloat(math.Abs(statement.Amount.Value), 'f', -1, 64)
		} else if statement.Currency != nil {
			creditCurrency = statement.Currency.Value
			creditAmount = strconv.FormatFloat(math.Abs(statement.Amount.Value), 'f', -1, 64)
			debitCurrency = creditCurrency
			debitAmount = creditAmount
		} else {
			creditCurrency = info.Currency
			creditAmount = strconv.FormatFloat(math.Abs(statement.Amount.Value), 'f', -1, 64)
			debitCurrency = creditCurrency
			debitAmount = creditAmount
		}

		idTransaction := info.IBAN + strconv.FormatInt(statement.TransactionID.Value, 10)

		if previousIDTransaction == "" {
			previousIDTransaction = idTransaction
		} else if previousIDTransaction != idTransaction {
			transaction := model.Transaction{
				Tenant:        workflow.Tenant,
				IDTransaction: previousIDTransaction,
				Transfers:     transfers,
			}
			log.Info().Msgf("Creating transaction %s in full %+v", transaction.IDTransaction, transaction)
			err := workflow.LedgerClient.CreateTransaction(transaction)
			if err != nil {
				log.Warn().Msgf("unable to create transaction %s/%s with error %s", workflow.Tenant, previousIDTransaction, err)
				return
			}
			workflow.Metrics.TransactionImported(len(transfers))

			for _, id := range idsNeedingConfirmation {
				err = workflow.PlaintextStorage.TouchFile("token/" + workflow.Token.ID + "/statements/" + id + "/ack_transfer")
				if err != nil {
					log.Warn().Err(err).Msgf("Unable to mark transfer discovery for %s/%s", workflow.Token.ID, id)
				}
			}
	
			previousIDTransaction = idTransaction
			transfers = make([]model.Transfer, 0)
			idsNeedingConfirmation = make([]string, 0)
		}

		if creditCurrency != debitCurrency {
			fmt.Printf("Currency mismatch credit %s %s debit %s %s\n", creditAmount, creditCurrency, debitAmount, debitCurrency)

			transfers = append(transfers, model.Transfer{
				IDTransfer: strconv.FormatInt(statement.TransferID.Value, 10) + "_1",
				Credit: model.AccountVault{
					Tenant: workflow.Tenant,
					Name:   debit,
				},
				Debit: model.AccountVault{
					Tenant: workflow.Tenant,
					Name:   info.BIC + "_" +  creditCurrency,
				},
				ValueDate: valueDate.Format("2006-01-02T15:04:05Z0700"),
				Amount:    creditAmount,
				Currency:  creditCurrency,
			})

			transfers = append(transfers, model.Transfer{
				IDTransfer: strconv.FormatInt(statement.TransferID.Value, 10) + "_2",
				Credit: model.AccountVault{
					Tenant: workflow.Tenant,
					Name:   info.BIC + "_" + debitCurrency,
				},
				Debit: model.AccountVault{
					Tenant: workflow.Tenant,
					Name:   credit,
				},
				ValueDate: valueDate.Format("2006-01-02T15:04:05Z0700"),
				Amount:    debitAmount,
				Currency:  debitCurrency,
			})

		} else {
			transfers = append(transfers, model.Transfer{
				IDTransfer: strconv.FormatInt(statement.TransferID.Value, 10) + "_1",
				Credit: model.AccountVault{
					Tenant: workflow.Tenant,
					Name:   credit,
				},
				Debit: model.AccountVault{
					Tenant: workflow.Tenant,
					Name:   debit,
				},
				ValueDate: valueDate.Format("2006-01-02T15:04:05Z0700"),
				Amount:    debitAmount,
				Currency:  debitCurrency,
			})
		}

		idsNeedingConfirmation = append(idsNeedingConfirmation, id)
	}

	if len(transfers) != 0 {
		transaction := model.Transaction{
			Tenant:        workflow.Tenant,
			IDTransaction: previousIDTransaction,
			Transfers:     transfers,
		}

		log.Info().Msgf("Creating transaction %s", transaction.IDTransaction)
		err := workflow.LedgerClient.CreateTransaction(transaction)
		if err != nil {
			log.Warn().Msgf("unable to create transaction %s/%s", workflow.Tenant, previousIDTransaction)
			return
		}
		workflow.Metrics.TransactionImported(len(transfers))
	}

	for _, id := range idsNeedingConfirmation {
		err = workflow.PlaintextStorage.TouchFile("token/" + workflow.Token.ID + "/statements/" + id + "/ack_transfer")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to mark transfer discovery for %s/%s", workflow.Token.ID, id)
		}
	}
}
