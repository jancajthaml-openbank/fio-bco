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
	Token 		     *model.Token
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
		Token: token,
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

	for _, transfer := range envelope.Statements {
		if transfer.TransferID == nil {
			continue
		}
		exists, err := workflow.PlaintextStorage.Exists("token/" + workflow.Token.ID + "/statements/" + envelope.IBAN + "/" + strconv.FormatInt(transfer.TransferID.Value, 10))
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to check if transaction %d exists for token %s IBAN %s", transfer.TransferID.Value, workflow.Token.ID, envelope.IBAN)
			return
		}
		if exists {
			continue
		}
		err = workflow.PlaintextStorage.TouchFile("token/" + workflow.Token.ID + "/statements/" + envelope.IBAN + "/" + strconv.FormatInt(transfer.TransferID.Value, 10) + "/mark")
		if err != nil {
			log.Warn().Err(err).Msgf("Unable to mark transaction %d as known for token %s IBAN %s", transfer.TransferID.Value, workflow.Token.ID, envelope.IBAN)
			return
		}
	}
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
