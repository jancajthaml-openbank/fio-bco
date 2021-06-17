// Copyright (c) 2016-2021, Jan Cajthaml <jan.cajthaml@gmail.com>
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

	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/support/http"
)

// Client represents fascade for vault http interactions
type Client struct {
	httpClient http.Client
	gateway    string
}

// NewClient returns new vault http client
func NewClient(gateway string) *Client {
	return &Client{
		gateway:    gateway,
		httpClient: http.NewClient(),
	}
}

// CreateAccount creates account in vault
func (client *Client) CreateAccount(account model.Account) error {
	if client == nil {
		return fmt.Errorf("nil deference")
	}

	req, err := http.NewRequest("POST", client.gateway+"/account/"+account.Tenant, account)
	if err != nil {
		return fmt.Errorf("create account error %w", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create account error %w", err)
	}

	if resp.StatusCode == 400 {
		return fmt.Errorf("create account malformed request")
	}
	if resp.StatusCode == 504 {
		return fmt.Errorf("create account timeout")
	}
	if resp.StatusCode != 200 && resp.StatusCode != 409 {
		return fmt.Errorf("create account invalid http status %s", resp.Status)
	}

	return nil
}
