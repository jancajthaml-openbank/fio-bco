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
)

// Configuration of application
type Configuration struct {
	// Tenant represent tenant of given vault
	Tenant string
	// SyncRate represents interval in which new statements are synchronised
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
	return loadConfFromEnv()
}
