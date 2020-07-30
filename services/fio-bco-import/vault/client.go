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

package vault

import (
	"fmt"

	"github.com/jancajthaml-openbank/fio-bco-import/http"
	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/utils"
)

// VaultClient represents fascade for http client
type VaultClient struct {
	underlying http.HttpClient
	gateway    string
}

// NewVaultClient returns new vault http client
func NewVaultClient(gateway string) VaultClient {
	return VaultClient{
		gateway:    gateway,
		underlying: http.NewHttpClient(),
	}
}

func (client VaultClient) CreateAccount(tenant string, account model.Account) error {
	request, err := utils.JSON.Marshal(account)
	if err != nil {
		return err
	}
	response, err := client.underlying.Post(client.gateway+"/account/"+tenant, request, nil)
	if err != nil {
		return fmt.Errorf("create account error %+v", err)
	}
	if response.Status == 400 {
		return fmt.Errorf("create account malformed request %+v", account)
	}
	if response.Status == 504 {
		return fmt.Errorf("create account timeout")
	}
	if response.Status != 200 && response.Status != 409 {
		return fmt.Errorf("create account error %s", response.String())
	}
	return nil
}
