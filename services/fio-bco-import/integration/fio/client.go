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

package fio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jancajthaml-openbank/fio-bco-import/model"
	"github.com/jancajthaml-openbank/fio-bco-import/support/http"
)

// Client represents fascade for FIO http interactions
type Client struct {
	gateway    string
	httpClient http.Client
}

// NewClient returns new fio http client
func NewClient(gateway string) *Client {
	return &Client{
		gateway:    gateway,
		httpClient: http.NewClient(),
	}
}

// GetStatementsEnvelope returns envelope since last synchronized id and sets that
// id as pivot for subsequent imports
func (client *Client) GetStatementsEnvelope(token model.Token) (*Envelope, error) {
	if client == nil {
		return nil, fmt.Errorf("nil deference")
	}

	var uri string
	if token.LastSyncedID != 0 {
		uri = fmt.Sprintf("/ib_api/rest/set-last-id/%s/%d/", token.Value, token.LastSyncedID)
	} else {
		uri = fmt.Sprintf("/ib_api/rest/set-last-date/%s/2012-07-27/", token.Value)
	}

	req, err := http.NewRequest("GET", client.gateway+uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if err == nil && resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("fio set last synced id invalid http status %s body %s", resp.Status, string(bodyBytes))
	}

	uri = "/ib_api/rest/last/" + token.Value + "/transactions.json"

	req, err = http.NewRequest("GET", client.gateway+uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err = client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 409 {
		return nil, fmt.Errorf("fio get transactions.json error token used before mandatory delay of 30 seconds")
	}

	if resp.StatusCode == 413 {
		return nil, fmt.Errorf("fio get transactions.json error more than 50k statements to be downloaded")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fio get transactions.json error invalid http status %s", resp.Status)
	}

	envelope := new(Envelope)
	err = json.NewDecoder(resp.Body).Decode(envelope)
	if err != nil {
		return nil, err
	}

	return envelope, nil
}
