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

func createAccountsFromStatements(
	tenant string,
	vaultClient *vault.Client,
	envelope *fio.Envelope,
) error {
	accounts := envelope.GetAccounts(tenant)

	for _, account := range accounts {
		log.Info().Msgf("Creating account %s", account.Name)

		err := vaultClient.CreateAccount(account)
		if err != nil {
			return fmt.Errorf("unable to create account %s with %w", account.Name, err)
		}
	}

	return nil
}

func createTransactionsFromStatements(
	tenant string,
	ledgerClient *ledger.Client,
	encryptedStorage localfs.Storage,
	metrics metrics.Metrics,
	token *model.Token,
	envelope *fio.Envelope,
) error {
	transactions := envelope.GetTransactions(tenant)

	for _, transaction := range transactions {
		log.Info().Msgf("Creating transaction %s", transaction.IDTransaction)
		err := ledgerClient.CreateTransaction(transaction)
		if err != nil {
			return fmt.Errorf("unable to create transaction %s/%s", transaction.Tenant, transaction.IDTransaction)
		}
		metrics.TransactionImported(len(transaction.Transfers))
		for _, transfer := range transaction.Transfers {
			if token.LastSyncedID > transfer.ID {
				continue
			}
			token.LastSyncedID = transfer.ID
			if !persistence.UpdateToken(encryptedStorage, token) {
				log.Warn().Msgf("unable to update token %s", token.ID)
			}
		}
	}

	return nil
}

// DownloadStatements download new statements from fio gateway
func (workflow Workflow) DownloadStatements() {
	if workflow.Token == nil {
		return
	}

	envelope, err := workflow.FioClient.GetStatementsEnvelope(*workflow.Token)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to get envelope")
		return
	}

	exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/info/" + envelope.Info.IBAN)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to check if info for token %s IBAN %s exists with error %s", workflow.Token.ID, envelope.Info.IBAN, err)
		return
	}

	if !exists {
		log.Info().Msgf("Should create info about synchronized account before downloading statements")
		data, err := json.Marshal(envelope.Info)
		if err != nil {
			log.Warn().Msgf("Unable to marshal info of %s/%s", workflow.Token.ID, envelope.Info.IBAN)
			return
		}
		err = workflow.PlaintextStorage.WriteFileExclusive("token/"+workflow.Token.ID+"/info/"+envelope.Info.IBAN, data)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to persist envelope info of %s/%s with error %s", workflow.Token.ID, envelope.Info.IBAN, err)
			return
		}
	}

	for _, transfer := range envelope.Statements {
		if transfer.TransferID == nil {
			continue
		}
		exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/statements/" + envelope.Info.IBAN + "/" + strconv.FormatInt(transfer.TransferID.Value, 10))
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to check if transaction %d exists for token %s IBAN %s with error %s", transfer.TransferID.Value, workflow.Token.ID, envelope.Info.IBAN, err)
			return
		}
		if exists {
			continue
		}
		data, err := json.Marshal(transfer)
		if err != nil {
			log.Warn().Msgf("Unable to marshal statement details of %s/%s/%d", workflow.Token.ID, envelope.Info.IBAN, transfer.TransferID.Value)
			continue
		}
		err = workflow.PlaintextStorage.WriteFileExclusive("token/"+workflow.Token.ID+"/statements/"+envelope.Info.IBAN+"/"+strconv.FormatInt(transfer.TransferID.Value, 10)+"/data", data)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to persist statement details of %s/%s/%d with error %s", workflow.Token.ID, envelope.Info.IBAN, transfer.TransferID.Value, err)
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
	//defer wg.Done()

	log.Debug().Msgf("token %s creating accounts from statements", workflow.Token.ID)

	IBANs, err := workflow.PlaintextStorage.ListDirectory("token/"+workflow.Token.ID+"/info", true)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to obtain IBANs from storage for token %s", workflow.Token.ID)
		return
	}

	for _, IBAN := range IBANs {

		data, err := workflow.PlaintextStorage.ReadFileFully("token/" + workflow.Token.ID + "/info/" + IBAN)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to load %s/%s info statement", workflow.Token.ID, IBAN)
			return
		}
		
		info := new(fio.Info)
		if json.Unmarshal(data, info) != nil {
			log.Warn().Msgf("Unable to unmarshal info %s/%s", workflow.Token.ID, IBAN)
			return
		}

		set := make(map[string]model.Account)
		idsNeedingConfirmation := make([]string, 0)

		ids, err := workflow.PlaintextStorage.ListDirectory("token/"+workflow.Token.ID+"/statements/"+IBAN, true)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to obtain transaction ids from storage for token %s IBAN %s", workflow.Token.ID, IBAN)
			return
		}

		// FIXME only once
		set[info.IBAN] = model.Account{
			Tenant:         workflow.Tenant,
			Name:           info.IBAN,
			Format:         "IBAN",
			Currency:       info.Currency,
			IsBalanceCheck: false,
		}

		for _, id := range ids {
			exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/statements/" + IBAN + "/" + id + "/accounts")
			if err != nil {
				log.Warn().Err(err).Msgf("Unable to check if statement %s/%s/%s accounts exists", workflow.Token.ID, IBAN, id)
				continue
			}
			if exists {
				continue
			}

			data, err := workflow.PlaintextStorage.ReadFileFully("token/" + workflow.Token.ID + "/statements/" + IBAN + "/" + id + "/data")
			if err != nil {
				log.Warn().Err(err).Msgf("Unable to load statement %s/%s/%s", workflow.Token.ID, IBAN, id)
				continue
			}

			statement := new(fio.Statement)
			if json.Unmarshal(data, statement) != nil {
				log.Warn().Msgf("Unable to unmarshal statement %s/%s/%s", workflow.Token.ID, IBAN, id)
				continue
			}

			log.Info().Msgf("Statement that accounts will be created from is %+v", statement)

			var normalizedAccount string
			var isIBAN bool
			var accountFormat string
			var currency string

			if statement.AccountTo == nil {
				// INFO fee and taxes and maybe card payments
				normalizedAccount = info.BIC
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

			if statement.Currency == nil {
				currency = info.Currency
			} else {
				currency = statement.Currency.Value
			}

			set[normalizedAccount] = model.Account{
				Tenant:         workflow.Tenant,
				Name:           normalizedAccount,
				Format:         accountFormat,
				Currency:       currency,
				IsBalanceCheck: false,
			}

			idsNeedingConfirmation = append(idsNeedingConfirmation, id)
		}

		for _, account := range set {
			log.Info().Msgf("Will create following account %+v", account)
		}

		if len(idsNeedingConfirmation) == 0 {
			continue
		}

		log.Info().Msgf("Following statements need confirmation %+v", idsNeedingConfirmation)

		for _, id := range idsNeedingConfirmation {
			err = workflow.PlaintextStorage.TouchFile("token/" + workflow.Token.ID + "/statements/" + IBAN + "/" + id + "/accounts")
			if err != nil {
				log.Warn().Err(err).Msgf("Unable to mark account discovery for %s/%s/%s", workflow.Token.ID, IBAN, id)
			}
		}
	}


	/*
	ids, err := plaintextStorage.ListDirectory("token/"+token.ID+"/statements/"+currency, true)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to obtain transaction ids from storage for token %s currency %s", token.ID, currency)
		return
	}

	accounts := make(map[string]bool)
	idsNeedingConfirmation := make([]string, 0)

	for _, id := range ids {
		exists, err := plaintextStorage.Exists("token/" + token.ID + "/statements/" + currency + "/" + id + "/accounts")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to check if statement %s/%s/%s accounts exists", token.ID, currency, id)
			continue
		}
		if exists {
			continue
		}

		data, err := plaintextStorage.ReadFileFully("token/" + token.ID + "/statements/" + currency + "/" + id + "/data")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to load statement %s/%s/%s", token.ID, currency, id)
			continue
		}

		statement := new(bondster.Statement)
		if json.Unmarshal(data, statement) != nil {
			log.Warn().Msgf("Unable to unmarshal statement %s/%s/%s", token.ID, currency, id)
			continue
		}

		accounts["BONDSTER_"+currency+"_TYPE_"+statement.Type] = true
		idsNeedingConfirmation = append(idsNeedingConfirmation, id)
	}

	if len(idsNeedingConfirmation) == 0 {
		return
	}

	accounts["BONDSTER_"+currency+"_TYPE_NOSTRO"] = true

	for account := range accounts {
		log.Info().Msgf("Creating account %s", account)

		request := model.Account{
			Tenant:         tenant,
			Name:           account,
			Currency:       currency,
			Format:         "BONDSTER_TECHNICAL",
			IsBalanceCheck: false,
		}
		err = vaultClient.CreateAccount(request)
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to create account %s/%s", tenant, account)
			return
		}
	}

	for _, id := range idsNeedingConfirmation {
		err = plaintextStorage.TouchFile("token/" + token.ID + "/statements/" + currency + "/" + id + "/accounts")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to mark account discovery for %s/%s/%s", token.ID, currency, id)
		}
	}
	*/
}

// SynchronizeStatements downloads new statements from fio gateway and creates accounts and transactions and normalizes them into value transfers
func (workflow Workflow) SynchronizeStatements() {
	if workflow.Token == nil {
		return
	}

	envelope, err := workflow.FioClient.GetStatementsEnvelope(*workflow.Token)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to get envelope")
		return
	}

	log.Debug().Msgf("token %s importing accounts", workflow.Token.ID)
	err = createAccountsFromStatements(workflow.Tenant, workflow.VaultClient, envelope)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to create accounts from envelope")
		return
	}

	log.Debug().Msgf("token %s importing transactions", workflow.Token.ID)
	err = createTransactionsFromStatements(workflow.Tenant, workflow.LedgerClient, workflow.EncryptedStorage, workflow.Metrics, workflow.Token, envelope)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to create transactions from envelope")
		return
	}
}
