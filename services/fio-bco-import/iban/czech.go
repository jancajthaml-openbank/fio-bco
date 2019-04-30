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

package iban

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var checkSumToString = make([]string, 99)

func init() {
	for i := 0; i < 10; i++ {
		checkSumToString[i] = "0" + strconv.Itoa(i)
	}

	for i := 10; i < 98; i++ {
		checkSumToString[i] = strconv.Itoa(i)
	}

	checkSumToString[98] = "98"
}

// CalculateCzech calculates IBAN for Czech Republic
func CalculateCzech(number, bankCode string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("iban.CalculateCzech(%s, %s) recovered in %+v", number, bankCode, r)
			result = ""
		}
	}()

	// canonise input
	canonisedNumber := strings.Replace(strings.Replace(number, "-", "", -1), " ", "", -1)
	// accountNumber of length 16
	paddedNumber := "0000000000000000"[0:16-len(canonisedNumber)] + canonisedNumber
	// bankCode of length 4
	paddedBankCode := "0000"[0:4-len(bankCode)] + bankCode
	// country code for "Czech Republic"
	countryCode := "CZ"
	// country code converted to digits
	countryDigits := "123500"
	// checksum mod 97
	checksum := (98 - mod97(paddedBankCode+paddedNumber+countryDigits))
	if checksum == 99 {
		return
	}

	result = countryCode + checkSumToString[checksum] + paddedBankCode + paddedNumber

	return
}
