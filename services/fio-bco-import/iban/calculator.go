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
		"0100", //KOMBCZPP - Komerční banka, a.s.
		"0300", //CEKOCZPP - Československá obchodní banka, a.s.
		"0600", //AGBACZPP - GE Money Bank, a.s.
		"0710", //CNBACZPP - Česká národní banka
		"0800", //GIBACZPX - Česká spořitelna, a.s.
		"2010", //FIOBCZPP - Fio, družstevní záložna
		"2020", //BOTKCZPP
		"2030", //?
		"2060", //CITFCZPP
		"2070", //MPUBCZPP
		"2100", //?
		"2200", //?
		"2220", //ARTTCZPP
		"2240", //POBNCZPP
		"2250", //CTASCZ22
		"2260", //?
		"2275", //?
		"2600", //CITICZPX - Citibank Europe plc, organizační složka
		"2700", //BACXCZPP - UniCredit Bank Czech Republic, a.s.
		"3030", //AIRACZPP - Air Bank. a.s.
		"3050", //BPPFCZP1
		"3060", //BPKOCZPP
		"3500", //INGBCZPP
		"4000", //EXPNCZPP
		"4300", //CMZRCZP1
		"5500", //RZBCCZPP - Raiffeisenbank, a.s.
		"5800", //JTBPCZPP
		"6000", //PMBPCZPP
		"6100", //EQBKCZPP
		"6200", //COBACZPX
		"6210", //BREXCZPP - mBank, a.s
		"6300", //GEBACZPP
		"6700", //SUBACZPP
		"6800", //VBOECZ2X
		"7910", //DEUTCZPX
		"7940", //SPWTCZ21
		"7950", //?
		"7960", //?
		"7970", //?
		"7980", //?
		"7990", //?
		"8030", //GENOCZ21
		"8040", //OBKLCZ2X
		"8060", //?
		"8090", //CZEECZPP
		"8150", //MIDLCZPP
		"8200", //?
		"8215", //?
		"8220", //PAERCZP1
		"8225", //ORRRCZP1
		"8230", //EEPSCZPP
		"8240", //?
		"8250", //BKCHCZPP
		"8260", //?
		"8265", //ICBKCZPP
		"8270", //FAPOCZP1
		"8280", //BEFKCZP1
		"8290", //ERSOCZPP
		"8291", //?
		"8292", //?
		"8293", //?
		"8294", //?
		"8295", //NVSRCZPP
		"8296", //?
		"8297", //?
		"8298", //ANCSCZP1
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
