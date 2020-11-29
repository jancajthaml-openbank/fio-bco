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

package integration

import (
	"time"

	"github.com/jancajthaml-openbank/fio-bco-import/persistence"

	localfs "github.com/jancajthaml-openbank/local-fs"
)

// FioImport represents fio gateway to ledger-rest import subroutine
type FioImport struct {
	callback   func(token string)
	storage    localfs.Storage
}

// NewFioImport returns fio import fascade
func NewFioImport(rootStorage string, storageKey []byte, callback func(token string)) *FioImport {
	storage, err := localfs.NewEncryptedStorage(rootStorage, storageKey)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	return &FioImport{
		callback: callback,
		storage:  storage,
	}
}

func (fio FioImport) getActiveTokens() ([]string, error) {
	tokens, err := persistence.LoadTokens(fio.storage)
	if err != nil {
		return nil, err
	}
	notBefore := time.Now().Add(time.Second * time.Duration(-6))
	uniq := make([]string, 0)
	visited := make(map[string]bool)
	for _, token := range tokens {
		if !token.CreatedAt.Before(notBefore) {
			continue
		}
		if _, ok := visited[token.Value]; !ok {
			visited[token.Value] = true
			uniq = append(uniq, token.ID)
		}
	}
	return uniq, nil
}

func (fio FioImport) importRoundtrip() {
	tokens, err := fio.getActiveTokens()
	if err != nil {
		log.Error().Msgf("unable to get active tokens %+v", err)
		return
	}

	for _, item := range tokens {
		log.Debug().Msgf("Request to import token %s", item)
		fio.callback(item)
	}
}

func (fio FioImport) Setup() error {
	return nil
}

func (fio FioImport) Work() {
	fio.importRoundtrip()
}

func (fio FioImport) Cancel() {

}

func (fio FioImport) Done() <-chan interface{} {
	done := make(chan interface{})
	close(done)
	return done
}
