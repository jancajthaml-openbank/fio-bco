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

//import (
	//"strconv"
//	"fmt"
	//"strings"
//)

// ValidateNetherlandsIBAN validates if number is IBAN
func ValidateNetherlandsIBAN(number string) bool {
	if len(number) != 18 {
		return false
	}
	if number[0:2] != "NL" {
		return false
	}
	switch number[4:8] {
	case
		"ABNA",
		"AFSR",
		"ABRF",
		"ACCI",
		"AGRV",
		"ACTO",
		"ACHL",
		"AHBK",
		"ARBN",
		"ACOA",
		"AVBG",
		"ADRB",
		"AEGO",
		"AEFD",
		"AEIM",
		"AENV",
		"AESP",
		"XAEX",
		"AFBR",
		"ACMN",
		"AFCR",
		"AFMO",
		"AGBS",
		"AKBK",
		"AKZO",
		"ASNN",
		"AETR",
		"AOIN",
		"ANAA",
		"ALPP",
		"ATRV",
		"AMRS",
		"ACMG",
		"AFAM",
		"AMSC",
		"AMBR",
		"XACE",
		"AITO",
		"STOL",
		"AMEF",
		"XAMS",
		"AMSG",
		"ANDL",
		"AAMM",
		"AAEM",
		"ANHH",
		"AOSE",
		"AOTS",
		"ABPT",
		"ARAM",
		"ARSN",
		"ASNE",
		"ASRB",
		"FTSI",
		"ASEG",
		"ASCE",
		"ATVE",
		"AVBD",
		"AZLH",
		"AZLV",
		"DCOM",
		"BKEM",
		"ESAO",
		"EXES",
		"BAGP",
		"BEOO",
		"INSI",
		"LABC",
		"BKMG",
		"BOFA",
		"BKCH",
		"BOFS",
		"TAIP",
		"BOTK",
		"OYEN",
		"TECT",
		"ZEEL",
		"ARTE",
		"BCDM",
		"BAYR",
		"BETP",
		"BBSC",
		"BEDO",
		"BEOU",
		"BTRD",
		"BITD",
		"BEFE",
		"BEOP",
		"BCMN",
		"BCSV",
		"BISU",
		"BICK",
		"BLNE",
		"BLSG",
		"BCMT",
		"PARB",
		"BOOE",
		"BOEM",
		"BORH",
		"BOTR",
		"BOSM",
		"BOWR",
		"BOXC",
		"BEFT",
		"BZWN":

		// TODO continue from https://www.thebankcodes.com/swift_code/swiftresults.php?searchstring=Netherlands&page=9

		{
			break
		}
	default:
		{
			return false
		}
	}

	return asciimod97(number[4:18] + number[0:4]) == 1
}

// CalculateNetherlands calculates IBAN for Netherlands
func CalculateNetherlands(number string, bankCode string) (result string) {
	defer func() {
		if recover() != nil {
			result = ""
		}
	}()
	/*
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
	checksum := (98 - mod97(paddedBankCode+paddedNumber+countryDigits))
	if checksum == 99 {
		return
	}
	result = countryCode + checkSumToString[checksum] + paddedBankCode + paddedNumber
	*/
	return
}
