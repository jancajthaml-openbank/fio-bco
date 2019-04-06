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

// Client represents fascade for http client
type Client struct {
	underlying *http.Client
}

// NewClient returns new http client
func NewClient() Client {
	return Client{
		underlying: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify:       true,
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

// Post performs http POST request for given url with given body
func (client Client) Post(url string, body []byte) (contents []byte, code int, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

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

		if err != nil {
			contents = nil
		} else {
			contents, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}()

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err = client.underlying.Do(req)
	if err != nil {
		return
	}

	code = resp.StatusCode
	return
}

// Get performs http GET request for given url
func (client Client) Get(url string) (contents []byte, code int, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

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

		if err != nil {
			contents = nil
		} else {
			contents, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}()

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json")

	resp, err = client.underlying.Do(req)
	if err != nil {
		return
	}

	code = resp.StatusCode
	return
}
