// Copyright (c) 2016-2019, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func loadConfFromEnv() Configuration {
	logLevel := strings.ToUpper(getEnvString("FIO_BCO_LOG_LEVEL", "DEBUG"))
	encryptionKey := getEnvString("FIO_BCO_ENCRYPTION_KEY", "")
	secrets := getEnvString("FIO_BCO_SECRETS", "")
	rootStorage := getEnvString("FIO_BCO_STORAGE", "/data")
	lakeHostname := getEnvString("FIO_BCO_LAKE_HOSTNAME", "")
	port := getEnvInteger("FIO_BCO_HTTP_PORT", 4000)

	if lakeHostname == "" || secrets == "" || rootStorage == "" || encryptionKey == "" {
		log.Fatal("missing required parameter to run")
	}

	serverCert, err := ioutil.ReadFile(secrets + "/domain.local.crt")
	if err != nil {
		log.Fatalf("unable to load certificate %s/domain.local.crt with error %+v", secrets, err)
	}

	serverKey, err := ioutil.ReadFile(secrets + "/domain.local.key")
	if err != nil {
		log.Fatalf("unable to load certificate %s/domain.local.key with error %+v", secrets, err)
	}

	keyData, err := ioutil.ReadFile(encryptionKey)
	if err != nil {
		log.Fatalf("unable to load encryption key from %s", encryptionKey)
	}

	storageKey, err := hex.DecodeString(string(keyData))
	if err != nil {
		log.Fatalf("invalid encryption key %+v at %s", err, encryptionKey)
	}

	return Configuration{
		RootStorage:   rootStorage,
		EncryptionKey: []byte(storageKey),
		ServerPort:    port,
		SecretKey:     serverKey,
		SecretCert:    serverCert,
		LakeHostname:  lakeHostname,
		LogLevel:      logLevel,
	}
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
		log.Panicf("invalid value of variable %s", key)
	}
	return cast
}
