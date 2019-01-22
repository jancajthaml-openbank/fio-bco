// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package utils

// TokensPath returns filepath of tokens for given tenant
func TokensPath(tenant string) string {
	return tenant + "/import/fio/token"
}

// TokenPath returns filepath of token for given tenant and token
func TokenPath(tenant, value string) string {
	return tenant + "/import/fio/token/" + value
}
