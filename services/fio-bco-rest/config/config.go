// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"github.com/jancajthaml-openbank/fio-bco-rest/support/env"
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
	// MinFreeDiskSpace respresents threshold for minimum disk free space to
	// be possible operating
	MinFreeDiskSpace uint64
	// MinFreeMemory respresents threshold for minimum available memory to
	// be possible operating
	MinFreeMemory uint64
}

// LoadConfig loads application configuration
func LoadConfig() Configuration {
	return Configuration{
		RootStorage:      env.String("FIO_BCO_STORAGE", "/data"),
		EncryptionKey:    env.HexFile("FIO_BCO_ENCRYPTION_KEY", nil),
		ServerPort:       env.Int("FIO_BCO_HTTP_PORT", 4000),
		ServerKey:        env.String("FIO_BCO_SERVER_KEY", ""),
		ServerCert:       env.String("FIO_BCO_SERVER_CERT", ""),
		LakeHostname:     env.String("FIO_BCO_LAKE_HOSTNAME", "127.0.0.1"),
		LogLevel:         strings.ToUpper(env.String("FIO_BCO_LOG_LEVEL", "INFO")),
		MinFreeDiskSpace: env.Uint64("VAULT_STORAGE_THRESHOLD", 0),
		MinFreeMemory:    env.Uint64("VAULT_MEMORY_THRESHOLD", 0),
	}
}
