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
	"context"
	"crypto/x509"
	"fmt"
	_http "net/http"
	"net/url"
	"regexp"
)

var redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

var schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)

// CheckRetry represent predicate if failed request should retry
type CheckRetry func(ctx context.Context, resp *_http.Response, err error) (bool, error)

// DefaultRetryPolicy represent default predicate if failed request should retry
func DefaultRetryPolicy(ctx context.Context, resp *_http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, v
			}
			if schemeErrorRe.MatchString(v.Error()) {
				return false, v
			}
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, v
			}
		}
		return true, nil
	}
	if resp.StatusCode == _http.StatusTooManyRequests {
		return true, nil
	}
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}
	return false, nil
}
