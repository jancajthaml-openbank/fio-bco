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

// ValidateGermanIBAN validates if number is IBAN
func ValidateGermanIBAN(number string) bool {
	if len(number) != 22 {
		return false
	}
	if number[0:2] != "DE" {
		return false
	}

	switch number[4:5] {
	case "1", "2", "3", "4", "5", "6", "7", "8":
		{
			break
		}
	default:
		{
			return false
		}
	}

	switch number[7:8] {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		{
			break
		}
	default:
		{
			return false
		}
	}

	return asciimod97(number[4:22] + number[0:4]) == 1
}
