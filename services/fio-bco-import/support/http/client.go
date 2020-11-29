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

package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const handshakeTimeout = 10 * time.Second
const dialTimeout = 30 * time.Second
const requestTimeout = 120 * time.Second

// Client represents fascade for http client
type Client struct {
	underlying *http.Client
}

// NewHTTPClient returns new http client
func NewHTTPClient() Client {
	return Client{
		underlying: &http.Client{
			Timeout: requestTimeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: dialTimeout,
				}).DialContext,
				TLSHandshakeTimeout: handshakeTimeout,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify:       false,
					MinVersion:               tls.VersionTLS12,
					MaxVersion:               tls.VersionTLS12,
					PreferServerCipherSuites: false,
					CurvePreferences: []tls.CurveID{
						tls.CurveP521,
						tls.CurveP384,
						tls.CurveP256,
					},
					CipherSuites: CipherSuites,
				},
			},
		},
	}
}

// Post performs http POST request
func (client *Client) Post(url string, body []byte, headers map[string]string) (response Response, err error) {
	response = Response{
		Status: 0,
		Data:   nil,
		Header: make(map[string]string),
	}

	if client == nil {
		return response, fmt.Errorf("cannot call methods on nil reference")
	}

	var req *http.Request
	var resp *http.Response

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime error %+v", r)
		}
		if err != nil && resp != nil {
			_, err = io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		} else if resp == nil && err != nil {
			err = fmt.Errorf("runtime error, no response")
		}

		if err == nil {
			response.Data, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}()
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.underlying.Do(req)
	if err != nil {
		return
	}
	for k, v := range resp.Header {
		response.Header[k] = v[len(v)-1]
	}
	response.Status = resp.StatusCode
	return
}

// Get performs http GET request
func (client *Client) Get(url string, headers map[string]string) (response Response, err error) {
	response = Response{
		Status: 0,
		Data:   nil,
		Header: make(map[string]string),
	}

	if client == nil {
		err = fmt.Errorf("cannot call methods on nil reference")
		return
	}

	var req *http.Request
	var resp *http.Response

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("runtime error %+v", r)
		}

		if err != nil && resp != nil {
			_, err = io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		} else if resp == nil && err != nil {
			err = fmt.Errorf("runtime error, no response %+v", err)
		}

		if err == nil {
			response.Data, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}()

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.underlying.Do(req)
	if err != nil {
		return
	}
	for k, v := range resp.Header {
		response.Header[k] = v[len(v)-1]
	}
	response.Status = resp.StatusCode
	return
}
