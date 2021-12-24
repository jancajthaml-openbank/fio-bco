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

// Calculate calculates IBAN given number and bank identity code
func Calculate(number string, identityCode string) string {
	// https://www.cnb.cz/miranda2/export/sites/www.cnb.cz/en/payment_systems/accounts_bank_codes/download/bank_codes_CR_128.pdf

	switch identityCode {
	case
		"0100",
		"0300",
		"0600",
		"0710",
		"0800",
		"2010",
		"2020",
		"2030",
		"2060",
		"2070",
		"2100",
		"2200",
		"2220",
		"2240",
		"2250",
		"2260",
		"2275",
		"2600",
		"2700",
		"3030",
		"3050",
		"3060",
		"3500",
		"4000",
		"4300",
		"5500",
		"5800",
		"6000",
		"6100",
		"6200",
		"6210",
		"6300",
		"6700",
		"6800",
		"7910",
		"7940",
		"7950",
		"7960",
		"7970",
		"7980",
		"7990",
		"8030",
		"8040",
		"8060",
		"8090",
		"8150",
		"8200",
		"8215",
		"8220",
		"8225",
		"8230",
		"8240",
		"8250",
		"8260",
		"8265",
		"8270",
		"8280",
		"8290",
		"8291",
		"8292",
		"8293",
		"8294",
		"8295",
		"8296",
		"8297",
		"8298",
		"0730":
		{
			return CalculateCzech(number, identityCode)
		}
	default:
		{
			return ""
		}
	}
}
