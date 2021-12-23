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
	"bytes"
	"io"
	_url "net/url"
	_http "net/http"
)

// Request wraps net/http request
type Request struct {
	body func() (io.Reader, error)
	*_http.Request
}

// NewRequest creates new http.Request
func NewRequest(method string, url string, data []byte) (*Request, error) {
	httpReq, err := _http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	bodyReader := func() (io.Reader, error) {
		return bytes.NewReader(data), nil
	}
	u, err := _url.Parse(url)
	if err == nil {
		httpReq.Host = u.Hostname()
		httpReq.Header["authority"] = []string{httpReq.Host}
	}
	if len(data) > 0 {
		httpReq.ContentLength = int64(len(data))
	}
	return &Request{bodyReader, httpReq}, nil
}

// SetHeader sets request header
func (request *Request) SetHeader(key string, value string) {
	if request == nil {
		return
	}
	request.Header[key] = []string{value}
}

