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
func Calculate(number string, bic string, bankCode string) string {

	if bic == "" {
		return CalculateCzech(number, bankCode)
	}

	switch bic {
	case "CSCHUS6S", "CSCHUS6SINT":
		{
			// TODO calculate US IBAN
			return ""	
		}
	case "FIOBCZPP", "FIOBCZPPXXX":
		{
			return CalculateCzech(number, "2010")
		}
	case "KOMBCZPP", "KOMBCZPPXXX":
		{
			return CalculateCzech(number, "0100")
		}
	case "CEKOCZPP", "CEKOCZPPXXX":
		{
			return CalculateCzech(number, "0300")
		}
	case "AGBACZPP", "AGBACZPPXXX":
		{
			return CalculateCzech(number, "0600")
		}
	case "CNBACZPP", "CNBACZPPXXX":
		{
			return CalculateCzech(number, "0710")
		}
	case "GIBACZPX", "GIBACZPXXXX":
		{
			return CalculateCzech(number, "0800")
		}
	case "BOTKCZPP", "BOTKCZPPXXX":
		{
			return CalculateCzech(number, "2020")
		}
	case "CITFCZPP", "CITFCZPPXXX":
		{
			return CalculateCzech(number, "2060")
		}
	case "MPUBCZPP", "MPUBCZPPXXX":
		{
			return CalculateCzech(number, "2070")
		}
	case "ARTTCZPP", "ARTTCZPPXXX":
		{
			return CalculateCzech(number, "2220")
		}
	case "POBNCZPP", "POBNCZPPXXX":
		{
			return CalculateCzech(number, "2240")
		}
	case "CTASCZ22", "CTASCZ22XXX":
		{
			return CalculateCzech(number, "2250")
		}

	case "ZUNOCZPP", "ZUNOCZPPXXX":
		{
			return CalculateCzech(number, "2310")
		}
	case "CITICZPX", "CITICZPXXXX":
		{
			return CalculateCzech(number, "2600")
		}
	case "BACXCZPP", "BACXCZPPXXX":
		{
			return CalculateCzech(number, "2700")
		}
	case "AIRACZPP", "AIRACZPPXXX":
		{
			return CalculateCzech(number, "3030")
		}
	case "INGBCZPP", "INGBCZPPXXX":
		{
			return CalculateCzech(number, "3500")
		}
	case "SOLACZPP", "SOLACZPPXXX":
		{
			return CalculateCzech(number, "4000")
		}
	case "CMZRCZP1", "CMZRCZP1XXX":
		{
			return CalculateCzech(number, "4300")
		}
	case "RZBCCZPP", "RZBCCZPPXXX":
		{
			return CalculateCzech(number, "5500")
		}
	case "JTBPCZPP", "JTBPCZPPXXX":
		{
			return CalculateCzech(number, "5800")
		}
	case "PMBPCZPP", "PMBPCZPPXXX":
		{
			return CalculateCzech(number, "6000")
		}
	case "EQBKCZPP", "EQBKCZPPXXX":
		{
			return CalculateCzech(number, "6100")
		}
	case "COBACZPX", "COBACZPXXXX":
		{
			return CalculateCzech(number, "6200")
		}
	case "BREXCZPP", "BREXCZPPXXX":
		{
			return CalculateCzech(number, "6210")
		}
	case "GEBACZPP", "GEBACZPPXXX":
		{
			return CalculateCzech(number, "6300")
		}
	case "SUBACZPP", "SUBACZPPXXX":
		{
			return CalculateCzech(number, "6700")
		}
	case "VBOECZ2X", "VBOECZ2XXXX":
		{
			return CalculateCzech(number, "6800")
		}
	case "DEUTCZPX", "DEUTCZPXXXX":
		{
			return CalculateCzech(number, "7910")
		}
	case "SPWTCZ21", "SPWTCZ21XXX":
		{
			return CalculateCzech(number, "7940")
		}
	case "GENOCZ21", "GENOCZ21XXX":
		{
			return CalculateCzech(number, "8030")
		}
	case "OBKLCZ2X", "OBKLCZ2XXXX":
		{
			return CalculateCzech(number, "8040")
		}
	case "CZEECZPP", "CZEECZPPXXX":
		{
			return CalculateCzech(number, "8090")
		}
	case "MIDLCZPP", "MIDLCZPPXXX":
		{
			return CalculateCzech(number, "8150")
		}		
	default:
		{
			return ""
		}
	}
}
