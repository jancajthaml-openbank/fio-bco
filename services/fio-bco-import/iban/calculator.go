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

// Calculate calculates IBAN given number and bank identity code
func Calculate(number, identityCode string) string {
	// https://www.cnb.cz/miranda2/export/sites/www.cnb.cz/en/payment_systems/accounts_bank_codes/download/bank_codes_CR_128.pdf

	switch identityCode {
	case
		"0100", // Komerční banka, a.s.
		"0300", // Československá obchodní banka, a.s.
		"2600", // Citibank Europe plc, organizační složka
		"3030", // Air Bank. a.s.
		"2700", // UniCredit Bank Czech Republic, a.s.
		"0600", // GE Money Bank, a.s.
		"0800", // Česká spořitelna, a.s.
		"5500", // Raiffeisenbank, a.s.
		"6210", // mBank, a.s
		"2010", // Fio, družstevní záložna
		"0710", // Česká národní banka
		"0730": // Česká národní banka - Clearing centre
		{
			return CalculateCzech(number, identityCode)
		}

	default:
		{
			return ""
		}
	}
}
