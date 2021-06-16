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

package http

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	_http "net/http"
	"time"
)

// Client represents fascade for http client
type Client struct {
	underlying *_http.Client
	checkRetry CheckRetry
	backoff    Backoff
}

// NewClient returns new http client
func NewClient() Client {
	return Client{
		underlying: &_http.Client{
			Transport: &_http.Transport{
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
					CipherSuites: []uint16{
						tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
						tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
						tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
						tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
						tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
						tls.TLS_RSA_WITH_AES_128_CBC_SHA,
						tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
						tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
						tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					},
				},
			},
		},
		checkRetry: DefaultRetryPolicy,
		backoff:    DefaultBackoff,
	}
}

// Do perform http.Request
func (client *Client) Do(req *Request) (*_http.Response, error) {
	if client == nil {
		return nil, fmt.Errorf("nil deference")
	}
	log.Debug().Str("url", req.URL.String()).Str("method", req.Method).Msg("performing request")
	var resp *_http.Response
	var attempt int
	var shouldRetry bool
	var doErr error
	var checkErr error

	for i := 0; ; i++ {
		attempt++
		var code int
		if req.body != nil {
			body, err := req.body()
			if err != nil {
				client.underlying.CloseIdleConnections()
				return resp, err
			}
			if c, ok := body.(io.ReadCloser); ok {
				req.Body = c
			} else {
				req.Body = ioutil.NopCloser(body)
			}
		}
		resp, doErr = client.underlying.Do(req.Request)
		if resp != nil {
			code = resp.StatusCode
		}
		shouldRetry, checkErr = client.checkRetry(req.Context(), resp, doErr)
		if doErr != nil {
			log.Error().Err(doErr).Str("method", req.Method).Str("url", req.URL.String()).Msgf("request failed")
		}
		if !shouldRetry {
			break
		}
		if doErr == nil {
			client.drainBody(resp.Body)
		}
		wait := client.backoff(time.Millisecond, time.Second, i, resp)
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		if code > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, code)
		}
		log.Debug().Str("request", desc).Str("timeout", wait.String()).Msgf("retrying request")
		select {
		case <-req.Context().Done():
			client.underlying.CloseIdleConnections()
			return nil, req.Context().Err()
		case <-time.After(wait):
		}
		httpreq := *req.Request
		req.Request = &httpreq
	}
	if doErr == nil && checkErr == nil && !shouldRetry {
		return resp, nil
	}
	defer client.underlying.CloseIdleConnections()
	err := doErr
	if checkErr != nil {
		err = checkErr
	}
	if resp != nil {
		client.drainBody(resp.Body)
	}
	if err == nil {
		return nil, fmt.Errorf("%s %s giving up after %d attempt(s)", req.Method, req.URL, attempt)
	}
	return nil, fmt.Errorf("%s %s giving up after %d attempt(s): %w", req.Method, req.URL, attempt, err)
}

func (client *Client) drainBody(body io.ReadCloser) {
	if client == nil {
		return
	}
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, int64(4096)))
	if err != nil {
		log.Error().Err(err).Msg("error reading response body")
	}
}
