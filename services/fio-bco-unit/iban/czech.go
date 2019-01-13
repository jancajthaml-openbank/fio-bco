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

package iban

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func CalculateCzech(number, bankCode string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("iban.CalculateCzech(%s, %s) recovered in %+v", number, bankCode, r)
			result = ""
		}
	}()

	// canonise input
	canonisedNumber := strings.Replace(number, "-", "", -1)
	// accountNumber of length 16
	paddedNumber := "0000000000000000"[0:16-len(canonisedNumber)] + canonisedNumber
	// bankCode of length 4
	paddedBankCode := "0000"[0:4-len(bankCode)] + bankCode
	// country code for "Czech Republic"
	countryCode := "CZ"
	// country code converted to digits
	countryDigits := "123500"
	// checksum of length 2
	var paddedChecksum string
	// checksum mod 97
	checksum := (98 - mod97(paddedBankCode+paddedNumber+countryDigits))
	switch checksum {
	case 99: // 98 - -1
		return
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		{
			paddedChecksum = "0" + strconv.Itoa(checksum)
			break
		}
	default:
		{
			paddedChecksum = strconv.Itoa(checksum)
			break
		}

	}

	result = countryCode + paddedChecksum + paddedBankCode + paddedNumber
}
