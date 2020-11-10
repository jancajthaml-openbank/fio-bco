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

package config

import (
	"time"
	"strings"
	"io/ioutil"
	"encoding/hex"
	"fmt"
)

// Configuration of application
type Configuration struct {
	// Tenant represent tenant of given vault
	Tenant string
	// SyncRate represents interval in which new statements are synchronized
	SyncRate time.Duration
	// FioGateway represent fio gateway uri
	FioGateway string
	// LedgerGateway represent ledger-rest gateway uri
	LedgerGateway string
	// VaultGateway represent vault-rest gateway uri
	VaultGateway string
	// RootStorage gives where to store journals
	RootStorage string
	// EncryptionKey represents current encryption key
	EncryptionKey []byte
	// LakeHostname represent hostname of openbank lake service
	LakeHostname string
	// LogOutput represents log output
	LogOutput string
	// LogLevel ignorecase log level
	LogLevel string
	// MetricsRefreshRate represents interval in which in memory metrics should be
	// persisted to disk
	MetricsRefreshRate time.Duration
	// MetricsOutput represents output file for metrics persistence
	MetricsOutput string
}

// GetConfig loads application configuration
func GetConfig() Configuration {
	logLevel := strings.ToUpper(envString("FIO_BCO_LOG_LEVEL", "INFO"))
	encryptionKey := envString("FIO_BCO_ENCRYPTION_KEY", "")
	rootStorage := envString("FIO_BCO_STORAGE", "/data")
	tenant := envString("FIO_BCO_TENANT", "")
	fioGateway := envString("FIO_BCO_FIO_GATEWAY", "https://www.fio.cz/ib_api/rest")
	ledgerGateway := envString("FIO_BCO_LEDGER_GATEWAY", "https://127.0.0.1:4401")
	vaultGateway := envString("FIO_BCO_VAULT_GATEWAY", "https://127.0.0.1:4400")
	syncRate := envDuration("FIO_BCO_SYNC_RATE", 22*time.Second)
	lakeHostname := envString("FIO_BCO_LAKE_HOSTNAME", "")
	metricsOutput := envFilename("FIO_BCO_METRICS_OUTPUT", "/tmp")
	metricsRefreshRate := envDuration("FIO_BCO_METRICS_REFRESHRATE", time.Second)

	if tenant == "" || lakeHostname == "" || rootStorage == "" || encryptionKey == "" {
		log.Error().Msg("missing required parameter to run")
		panic("missing required parameter to run")
	}

	keyData, err := ioutil.ReadFile(encryptionKey)
	if err != nil {
		log.Error().Msgf("unable to load encryption key from %s", encryptionKey)
		panic(fmt.Sprintf("unable to load encryption key from %s", encryptionKey))
	}

	key, err := hex.DecodeString(string(keyData))
	if err != nil {
		log.Error().Msgf("invalid encryption key %+v at %s", err, encryptionKey)
		panic(fmt.Sprintf("invalid encryption key %+v at %s", err, encryptionKey))
	}

	return Configuration{
		Tenant:             tenant,
		RootStorage:        rootStorage + "/t_" + tenant + "/import/fio",
		EncryptionKey:      []byte(key),
		FioGateway:         fioGateway,
		SyncRate:           syncRate,
		LedgerGateway:      ledgerGateway,
		VaultGateway:       vaultGateway,
		LakeHostname:       lakeHostname,
		LogLevel:           logLevel,
		MetricsRefreshRate: metricsRefreshRate,
		MetricsOutput:      metricsOutput,
	}
}
