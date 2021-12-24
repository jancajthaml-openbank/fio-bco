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

package iban

import "strings"

// ValidateCzechIBAN validates if number is IBAN
func ValidateCzechIBAN(number string) bool {
	if len(number) != 24 {
		return false
	}
	if number[0:2] != "CZ" {
		return false
	}

	switch number[4:5] {
	case "0", "2", "3", "4", "5", "7", "8":
		{
			break
		}
	default:
		{
			return false
		}
	}

	return asciimod97(number[4:24] + number[0:4]) == 1
}

// CalculateCzech calculates IBAN for Czech Republic
func CalculateCzech(number string, bankCode string) (result string) {
	defer func() {
		if recover() != nil {
			result = ""
		}
	}()
	// canonize input
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
	checksum := (98 - asciimod97(paddedBankCode+paddedNumber+countryDigits))
	if checksum == 99 {
		return
	}
	result = countryCode + checkSumToString[checksum] + paddedBankCode + paddedNumber
	return
}
