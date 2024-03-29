// Copyright (c) 2016-2023, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package actor

const (
	// ReqSynchronizeToken fio message request code for "Synchronize Token"
	ReqSynchronizeToken = "ST"
	// RespSynchronizeToken fio message response code for "Synchronize Token"
	RespSynchronizeToken = "TS"
	// ReqCreateToken fio message request code for "New Token"
	ReqCreateToken = "NT"
	// RespCreateToken fio message response code for "New Token"
	RespCreateToken = "TN"
	// ReqDeleteToken fio message request code for "Delete Token"
	ReqDeleteToken = "DT"
	// RespDeleteToken fio message response code for "Delete Token"
	RespDeleteToken = "TD"
	// RespTokenDoesNotExist fio message response code for "Token does not Exist"
	RespTokenDoesNotExist = "EM"
	// FatalError fio message response code for "Error"
	FatalError = "EE"
)
