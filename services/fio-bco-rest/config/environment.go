// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
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
	"os"
	"strconv"
	"strings"

	storage "github.com/jancajthaml-openbank/local-fs"
	log "github.com/sirupsen/logrus"
)

func loadConfFromEnv() Configuration {
	logOutput := getEnvString("FIO_BCO_LOG", "")
	logLevel := strings.ToUpper(getEnvString("FIO_BCO_LOG_LEVEL", "DEBUG"))
	secrets := getEnvString("FIO_BCO_SECRETS", "")
	rootStorage := getEnvString("FIO_BCO_STORAGE", "/data")
	lakeHostname := getEnvString("FIO_BCO_LAKE_HOSTNAME", "")
	port := getEnvInteger("FIO_BCO_HTTP_PORT", 443)

	if lakeHostname == "" || secrets == "" || rootStorage == "" {
		log.Fatal("missing required parameter to run")
	}

	if os.MkdirAll(rootStorage, os.ModePerm) != nil {
		log.Fatal("unable to assert storage directory")
	}

	cert, err := storage.ReadFileFully(secrets + "/domain.local.crt")
	if err != nil {
		log.Fatalf("unable to load certificate %s/domain.local.crt", secrets)
	}

	key, err := storage.ReadFileFully(secrets + "/domain.local.key")
	if err != nil {
		log.Fatalf("unable to load certificate %s/domain.local.key", secrets)
	}

	return Configuration{
		RootStorage:  rootStorage,
		ServerPort:   port,
		SecretKey:    key,
		SecretCert:   cert,
		LakeHostname: lakeHostname,
		LogOutput:    logOutput,
		LogLevel:     logLevel,
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
