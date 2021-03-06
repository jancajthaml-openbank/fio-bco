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
			if token.LastSyncedID > transfer.IDTransfer {
				continue
			}
			token.LastSyncedID = transfer.IDTransfer
			if !persistence.UpdateToken(encryptedStorage, token) {
				log.Warn().Msgf("unable to update token %s", token.ID)
			}
		}
	}

	return nil
}

func synchronizeNewStatements(
	encryptedStorage localfs.Storage,
	plaintextStorage localfs.Storage,
	token *model.Token,
	tenant string,
	fioClient *fio.Client,
	vaultClient *vault.Client,
	ledgerClient *ledger.Client,
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
