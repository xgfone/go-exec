// Copyright 2021 xgfone
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
	"strings"
	"testing"
)

const combinedOutputScript1 = `
ls README.md
ls notexistfile
ls notexitsfile 2>/dev/null
`

const combinedOutputScript2 = `
ls notexitsfile 2>/dev/null
ls notexistfile
ls README.md
`

func TestCombinedOutputCmdHook(t *testing.T) {
	cmd := DefaultCmd
	cmd.CmdHook = CombinedOutputCmdHook
	cmd.ResultHook = func(Result) {}

	stdout, err := cmd.OutputShellScript(context.TODO(), combinedOutputScript1)
	if err == nil {
		t.Error("not expect error is nil")
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) != 2 {
		t.Errorf("expect %d line, but got %d", 2, len(lines))
	}
	for i, line := range lines {
		switch i {
		case 0:
			if line != "README.md" {
				t.Errorf("expect line '%s', but got '%s'", "README.md", line)
			}
		case 1:
			if line == "notexistfile" {
				t.Errorf("unexpect line '%s'", "notexistfile")
			}
		}
	}

	stdout, err = cmd.OutputShellScript(context.TODO(), combinedOutputScript2)
	if err != nil {
		t.Errorf("expect error is nil, but got '%v'", err)
	}
	lines = strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) != 2 {
		t.Errorf("expect %d line, but got %d", 2, len(lines))
	}
	for i, line := range lines {
		switch i {
		case 0:
			if line == "notexistfile" {
				t.Errorf("unexpect line '%s'", "notexistfile")
			}
		case 1:
			if line != "README.md" {
				t.Errorf("expect line '%s', but got '%s'", "README.md", line)
			}
		}
	}
}
