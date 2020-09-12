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
	"encoding/hex"
	"io/ioutil"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func loadConfFromEnv() Configuration {
	logLevel := strings.ToUpper(getEnvString("FIO_BCO_LOG_LEVEL", "DEBUG"))
	encryptionKey := getEnvString("FIO_BCO_ENCRYPTION_KEY", "")
	serverKey := getEnvString("FIO_BCO_SERVER_KEY", "")
	serverCert := getEnvString("FIO_BCO_SERVER_CERT", "")
	rootStorage := getEnvString("FIO_BCO_STORAGE", "/data")
	lakeHostname := getEnvString("FIO_BCO_LAKE_HOSTNAME", "")
	port := getEnvInteger("FIO_BCO_HTTP_PORT", 4000)
	minFreeDiskSpace := getEnvInteger("VAULT_STORAGE_THRESHOLD", 0)
	minFreeMemory := getEnvInteger("VAULT_MEMORY_THRESHOLD", 0)
	metricsOutput := getEnvFilename("FIO_BCO_METRICS_OUTPUT", "/tmp")
	metricsRefreshRate := getEnvDuration("FIO_BCO_METRICS_REFRESHRATE", time.Second)

	if lakeHostname == "" || serverKey == "" || serverCert == "" || rootStorage == "" || encryptionKey == "" {
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
		RootStorage:        rootStorage,
		EncryptionKey:      []byte(key),
		ServerPort:         port,
		ServerKey:          serverKey,
		ServerCert:         serverCert,
		LakeHostname:       lakeHostname,
		LogLevel:           logLevel,
		MetricsRefreshRate: metricsRefreshRate,
		MetricsOutput:      metricsOutput,
		MinFreeDiskSpace:   uint64(minFreeDiskSpace),
		MinFreeMemory:      uint64(minFreeMemory),
	}
}

func getEnvFilename(key, fallback string) string {
	var value = strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	value = filepath.Clean(value)
	if os.MkdirAll(value, os.ModePerm) != nil {
		return fallback
	}
	return value
}

func getEnvString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInteger(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	cast, err := strconv.Atoi(value)
	if err != nil {
		log.Error().Msgf("invalid value of variable %s", key)
		return fallback
	}
	return cast
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	cast, err := time.ParseDuration(value)
	if err != nil {
		log.Error().Msgf("invalid value of variable %s", key)
		return fallback
	}
	return cast
}
