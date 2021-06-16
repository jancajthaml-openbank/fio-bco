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
	"math"
	_http "net/http"
	"strconv"
	"time"
)

// Backoff represents http request backoff policy
type Backoff func(min time.Duration, max time.Duration, attemptNum int, resp *_http.Response) time.Duration

// DefaultBackoff represents default http request backoff policy
func DefaultBackoff(min time.Duration, max time.Duration, attemptNum int, resp *_http.Response) time.Duration {
	if resp != nil {
		if resp.StatusCode == _http.StatusTooManyRequests || resp.StatusCode == _http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return time.Second * time.Duration(sleep)
				}
			}
		}
	}
	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}
