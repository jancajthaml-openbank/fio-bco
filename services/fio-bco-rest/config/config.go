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
)

// Configuration of application
type Configuration struct {
	// RootStorage gives where to store journals
	RootStorage string
	// EncryptionKey represents current encryption key
	EncryptionKey []byte
	// ServerPort is port which server is bound to
	ServerPort int
	// ServerKey path to server tls key file
	ServerKey string
	// ServerCert path to server tls cert file
	ServerCert string
	// LakeHostname represent hostname of openbank lake service
	LakeHostname string
	// LogLevel ignorecase log level
	LogLevel string
	// MetricsRefreshRate represents interval in which in memory metrics should be
	// persisted to disk
	MetricsRefreshRate time.Duration
	// MetricsOutput represents output file for metrics persistence
	MetricsOutput string
	// MinFreeDiskSpace respresents threshold for minimum disk free space to
	// be possible operating
	MinFreeDiskSpace uint64
	// MinFreeMemory respresents threshold for minimum available memory to
	// be possible operating
	MinFreeMemory uint64
}

// GetConfig loads application configuration
func GetConfig() Configuration {
	return Configuration{
		RootStorage:        envString("FIO_BCO_STORAGE", "/data"),
		EncryptionKey:      envSecret("FIO_BCO_ENCRYPTION_KEY", nil),
		ServerPort:         envInteger("FIO_BCO_HTTP_PORT", 4000),
		ServerKey:          envString("FIO_BCO_SERVER_KEY", ""),
		ServerCert:         envString("FIO_BCO_SERVER_CERT", ""),
		LakeHostname:       envString("FIO_BCO_LAKE_HOSTNAME", "127.0.0.1"),
		LogLevel:           strings.ToUpper(envString("FIO_BCO_LOG_LEVEL", "INFO")),
		MetricsRefreshRate: envDuration("FIO_BCO_METRICS_REFRESHRATE", time.Second),
		MetricsOutput:      envFilename("FIO_BCO_METRICS_OUTPUT", "/tmp/fio-bco-rest-metrics"),
		MinFreeDiskSpace:   uint64(envInteger("VAULT_STORAGE_THRESHOLD", 0)),
		MinFreeMemory:      uint64(envInteger("VAULT_MEMORY_THRESHOLD", 0)),
	}
}
