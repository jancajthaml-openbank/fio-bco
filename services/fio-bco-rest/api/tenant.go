// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/jancajthaml-openbank/fio-bco-rest/system"

	"github.com/labstack/echo/v4"
)

// CreateTenant enables fio-bco-import@{tenant}
func CreateTenant(control system.Control) func(c echo.Context) error {
	if control == nil {
		return func(c echo.Context) error {
			c.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
	}
	return func(c echo.Context) error {
		unescapedTenant, err := url.PathUnescape(c.Param("tenant"))
		if err != nil {
			c.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
		tenant := strings.TrimSpace(unescapedTenant)
		if tenant == "" {
			c.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
		err = control.EnableUnit("import@" + tenant + ".service")
		if err != nil {
			return err
		}
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}
}

// DeleteTenant disables fio-bco-import@{tenant}
func DeleteTenant(control system.Control) func(c echo.Context) error {
	if control == nil {
		return func(c echo.Context) error {
			c.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
	}
	return func(c echo.Context) error {
		unescapedTenant, err := url.PathUnescape(c.Param("tenant"))
		if err != nil {
			c.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
		tenant := strings.TrimSpace(unescapedTenant)
		if tenant == "" {
			c.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
		err = control.DisableUnit("import@" + tenant + ".service")
		if err != nil {
			return err
		}
		c.Response().WriteHeader(http.StatusOK)
		return nil
	}
}

// ListTenants lists fio-bco-import@
func ListTenants(control system.Control) func(c echo.Context) error {
	if control == nil {
		return func(c echo.Context) error {
			c.Response().WriteHeader(http.StatusNotFound)
			return nil
		}
	}
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		units, err := control.ListUnits("import@")
		if err != nil {
			return err
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		for idx, unit := range units {
			if idx == len(units)-1 {
				c.Response().Write([]byte(unit))
			} else {
				c.Response().Write([]byte(unit))
				c.Response().Write([]byte("\n"))
			}
			c.Response().Flush()
		}

		return nil
	}
}
