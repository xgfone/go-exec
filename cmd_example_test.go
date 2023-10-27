// Copyright 2021~2023 xgfone
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

//go:build unix || aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build unix aix darwin dragonfly freebsd linux netbsd openbsd solaris

package exec

import (
	"context"
	"fmt"
)

const scripttmpl = `
ls %s
rm -rf %s
`

func ExampleCmd() {
	_ = Execute(context.Background(), "mkdir", "testdir")
	_ = ExecuteShellCmd(context.Background(), "echo abc > %s/%s", "testdir", "testfile")

	data, _ := OutputShellCmd(context.Background(), "cat %s/%s", "testdir", "testfile")
	fmt.Println(data)

	_, _, err := RunShellScript(context.Background(), scripttmpl, "testdir", "testdir")
	fmt.Println(err)

	// Output:
	// abc
	//
	// <nil>
}
