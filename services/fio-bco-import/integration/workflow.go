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
	"sort"

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
		FioClient:        http.NewFioClient(fioGateway, *token),
		VaultClient:      http.NewVaultClient(vaultGateway),
		LedgerClient:     http.NewLedgerClient(ledgerGateway),
		EncryptedStorage: encryptedStorage,
		PlaintextStorage: plaintextStorage,
		Metrics:          metrics,
	}
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

	statements, err := fioClient.GetTransactions()
	if err != nil {
		return
	}
	if len(statements.Statement.TransactionList.Transactions) == 0 {
		return
	}

	log.Debug().Msgf("token %s sorting statements", token.ID)

	sort.SliceStable(statements.Statement.TransactionList.Transactions, func(i, j int) bool {
		return statements.Statement.TransactionList.Transactions[i].TransactionID.Value < statements.Statement.TransactionList.Transactions[j].TransactionID.Value
	})

	log.Debug().Msgf("token %s importing accounts", token.ID)

	for account := range statements.GetAccounts(tenant) {
		err = vaultClient.CreateAccount(account)
		if err != nil {
			log.Warn().Msgf("Unable to create account %s/%s with %+v", tenant, account.Name, err)
			return
		}
	}

	log.Debug().Msgf("token %s importing transactions", token.ID)

	for transaction := range statements.GetTransactions(tenant) {
		err = ledgerClient.CreateTransaction(transaction)
		if err != nil {
			log.Warn().Msgf("Unable to create transaction %s/%s with %+v", tenant, transaction.IDTransaction, err)
			return
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

}

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
