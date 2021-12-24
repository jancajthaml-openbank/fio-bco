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

import "strconv"

var checkSumToString = make([]string, 99)

func init() {
	for i := 0; i < 10; i++ {
		checkSumToString[i] = "0" + strconv.Itoa(i)
	}
	for i := 10; i < 99; i++ {
		checkSumToString[i] = strconv.Itoa(i)
	}
}

func asciimod97(number string) int {
	var (
		d uint
		i int
		x uint
		l = len(number)
	)

scan:
	d = uint(number[i]) - 48
	if d > 9 {
		x = (((x * 100) + d - 7) % 97)
	} else {
		x = (((x * 10) + d) % 97)
	}
	i++
	if i != l {
		goto scan
	}
	return int(x)
}