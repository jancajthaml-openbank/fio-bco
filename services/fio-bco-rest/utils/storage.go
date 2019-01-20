// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
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

package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"syscall"
	"unsafe"
)

var defaultBufferSize int

func init() {
	defaultBufferSize = 2 * os.Getpagesize()
}

func nameFromDirent(dirent *syscall.Dirent) []byte {
	reg := int(uint64(dirent.Reclen) - uint64(unsafe.Offsetof(syscall.Dirent{}.Name)))

	var name []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&name))
	header.Cap = reg
	header.Len = reg
	header.Data = uintptr(unsafe.Pointer(&dirent.Name[0]))

	if index := bytes.IndexByte(name, 0); index >= 0 {
		header.Cap = index
		header.Len = index
	}

	return name
}

// ListDirectory returns sorted slice of item names in given absolute path
// default sorting is ascending
func ListDirectory(absPath string, ascending bool) (result []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = fmt.Errorf("ListDirectory paniced: %+v", r)
			}
		}
		if err != nil {
			result = nil
		}
	}()

	var (
		n  int
		dh *os.File
		de *syscall.Dirent
	)

	dh, err = os.Open(filepath.Clean(absPath))
	if err != nil {
		return
	}

	fd := int(dh.Fd())
	result = make([]string, 0)

	scratchBuffer := make([]byte, defaultBufferSize)

	for {
		n, err = syscall.ReadDirent(fd, scratchBuffer)
		if err != nil {
			if r := dh.Close(); r != nil {
				err = r
			}
			return
		}
		if n <= 0 {
			break
		}
		buf := scratchBuffer[:n]
		for len(buf) > 0 {
			de = (*syscall.Dirent)(unsafe.Pointer(&buf[0]))
			buf = buf[de.Reclen:]

			if de.Ino == 0 {
				continue
			}

			nameSlice := nameFromDirent(de)
			switch len(nameSlice) {
			case 0:
				continue
			case 1:
				if nameSlice[0] == '.' {
					continue
				}
			case 2:
				if nameSlice[0] == '.' && nameSlice[1] == '.' {
					continue
				}
			}
			result = append(result, string(nameSlice))
		}
	}

	if r := dh.Close(); r != nil {
		err = r
		return
	}

	if ascending {
		sort.Slice(result, func(i, j int) bool {
			return result[i] < result[j]
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return result[i] > result[j]
		})
	}

	return
}

// ReadFileFully reads whole file given absolute path
func ReadFileFully(absPath string) ([]byte, error) {
	f, err := os.OpenFile(absPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, fi.Size())
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf, nil
}

// Exists returns true if absolute path exists
func Exists(absPath string) bool {
	_, err := os.Stat(absPath)
	return !os.IsNotExist(err)
}

// NotExists returns true if absolute path does not exists
func NotExists(absPath string) bool {
	_, err := os.Stat(absPath)
	return os.IsNotExist(err)
}
