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

package actor

import (
	"github.com/jancajthaml-openbank/fio-bco-import/metrics"

	system "github.com/jancajthaml-openbank/actor-system"
	localfs "github.com/jancajthaml-openbank/local-fs"
)

// System represents actor system subroutine
type System struct {
	system.System
	Tenant        string
	Storage       localfs.Storage
	Metrics       metrics.Metrics
	FioGateway    string
	LedgerGateway string
	VaultGateway  string
}

// NewActorSystem returns actor system fascade
func NewActorSystem(tenant string, lakeEndpoint string, fioEndpoint string, vaultEndpoint string, ledgerEndpoint string, rootStorage string, storageKey []byte, metrics metrics.Metrics) *System {
	storage, err := localfs.NewPlaintextStorage(rootStorage)
	if err != nil {
		log.Error().Msgf("Failed to ensure storage %+v", err)
		return nil
	}
	sys, err := system.New("FioImport/"+tenant, lakeEndpoint)
	if err != nil {
		log.Error().Msgf("Failed to register actor system %+v", err)
		return nil
	}
	result := new(System)
	result.System = sys
	result.Storage = storage
	result.Metrics = metrics
	result.Tenant = tenant
	result.FioGateway = fioEndpoint
	result.LedgerGateway = ledgerEndpoint
	result.VaultGateway = vaultEndpoint
	result.System.RegisterOnMessage(ProcessMessage(result))
	return result
}

// Setup does nothing
func (system *System) Setup() error {
	return nil
}

// Work starts actor system
func (system *System) Work() {
	if system == nil {
		return
	}
	system.System.Start()
}

// Cancel does nothing
func (system *System) Cancel() {
	if system == nil {
		return
	}
	system.System.Stop()
}

// Done always returns done
func (system *System) Done() <-chan interface{} {
	done := make(chan interface{})
	close(done)
	return done
}
