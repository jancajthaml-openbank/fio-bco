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

package fio

import (
  "fmt"

  "github.com/jancajthaml-openbank/fio-bco-import/http"
  "github.com/jancajthaml-openbank/fio-bco-import/model"
  "github.com/jancajthaml-openbank/fio-bco-import/utils"
)

// FioClient represents fascade for http client
type FioClient struct {
  underlying http.HttpClient
  gateway    string
  token      model.Token
}

// NewFioClient returns new bondster http client
func NewFioClient(gateway string, token model.Token) FioClient {
  return FioClient{
    gateway:    gateway,
    underlying: http.NewHttpClient(),
    token:      token,
  }
}

func (client *FioClient) setLastSyncedID() error {
  if client == nil {
    return fmt.Errorf("nil deference")
  }

  var uri string
  if client.token.LastSyncedID != 0 {
    uri = fmt.Sprintf("/ib_api/rest/set-last-id/%s/%d/", client.token.Value, client.token.LastSyncedID)
  } else {
    uri = fmt.Sprintf("/ib_api/rest/set-last-date/%s/2012-07-27/", client.token.Value)
  }

  response, err := client.underlying.Get(client.gateway+uri, nil)
  if err != nil {
    return err
  }
  if response.Status != 200 {
    return fmt.Errorf("fio set last synced id error %s", response.String())
  }
  return nil
}

func (client *FioClient) GetTransactions() (*FioImportEnvelope, error) {
  if client == nil {
    return nil, fmt.Errorf("nil deference")
  }

  err := client.setLastSyncedID()
  if err != nil {
    return nil, err
  }

  response, err := client.underlying.Get(client.gateway+"/ib_api/rest/last/" + client.token.Value + "/transactions.json", nil)
  if err != nil {
    return nil, err
  }
  if response.Status != 200 {
    return nil, fmt.Errorf("fio set last synced id error %s", response.String())
  }

  var envelope = new(FioImportEnvelope)
  err = utils.JSON.Unmarshal(response.Data, envelope)
  if err != nil {
    return nil, err
  }

  return envelope, nil
}
