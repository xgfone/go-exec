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

package exec

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// DefaultShell is the default shell to execute the shell command or script.
var DefaultShell = "bash"

// ShellScriptDir is the directory to save the script file to be executed.
//
// If OS is windows or js, it is reset to "". But you can set it to somewhere.
var ShellScriptDir = os.TempDir()

func init() {
	switch runtime.GOOS {
	case "js", "windows":
		ShellScriptDir = ""
	}
}

// RunShellCmd runs the command cmdargs with cmdargs as the shell command,
// that's,
//   shell -c "fmt.Sprintf(cmdfmt, cmdargs...)".
func (c Cmd) RunShellCmd(ctx context.Context, cmdfmt string, cmdargs ...string) (
	stdout, stderr string, err error) {
	shell := c.Shell
	if shell == "" {
		shell = DefaultShell
	}

	if _len := len(cmdargs); _len != 0 {
		vs := make([]interface{}, _len)
		for i := 0; i < _len; i++ {
			vs[i] = cmdargs[i]
		}
		cmdfmt = fmt.Sprintf(cmdfmt, vs...)
	}

	return c.Run(ctx, shell, "-c", cmdfmt)
}

// ExecuteShellCmd is the same as RunShellCmd, but only returns the error.
func (c Cmd) ExecuteShellCmd(cxt context.Context, cmdfmt string, cmdargs ...string) error {
	_, _, err := c.RunShellCmd(cxt, cmdfmt, cmdargs...)
	return err
}

// OutputShellCmd is the same as RunShellCmd, but only returns stdout and error.
func (c Cmd) OutputShellCmd(cxt context.Context, cmdfmt string, cmdargs ...string) (string, error) {
	stdout, _, err := c.RunShellCmd(cxt, cmdfmt, cmdargs...)
	return string(stdout), err
}

// ExecuteShellScript is the same as RunShellScript, but only returns the error.
func (c Cmd) ExecuteShellScript(ctx context.Context, script string, args ...string) error {
	_, _, err := c.RunShellScript(ctx, script, args...)
	return err
}

// OutputShellScript is the same as RunShellScript, but only returns stdout and error.
func (c Cmd) OutputShellScript(ctx context.Context, script string, args ...string) (
	stdout string, err error) {
	stdout, _, err = c.RunShellScript(ctx, script, args...)
	return
}

// RunShellScript runs the script with args as the shell script,
// the content of which is fmt.Sprintf(script, args...).
func (c Cmd) RunShellScript(ctx context.Context, script string, args ...string) (
	stdout, stderr string, err error) {
	_script := script
	if _len := len(args); _len != 0 {
		vs := make([]interface{}, _len)
		for i := 0; i < _len; i++ {
			vs[i] = args[i]
		}
		_script = fmt.Sprintf(script, vs...)
	}

	filename, err := c.getScriptFile(_script)
	if err != nil {
		err = NewResult(script, args, "", "", err)
		return
	}
	defer os.RemoveAll(filename)

	shell := c.Shell
	if shell == "" {
		shell = DefaultShell
	}
	return c.Run(ctx, shell, filename)
}

func (c Cmd) getScriptFile(script string) (filename string, err error) {
	var buf [8]byte
	const chars = `0123456789abcdefghijklmnopqrstuvwxyz`
	const charslen = len(chars)
	for i := 0; i < 8; i++ {
		buf[i] = chars[rand.Intn(charslen)]
	}

	data := []byte(script)
	md5sum := md5.Sum(data)
	hexsum := hex.EncodeToString(md5sum[:])
	filename = fmt.Sprintf("_cmd_exec_run_shell_script_md5%s_%d_%s.sh",
		hexsum, time.Now().UnixMicro(), string(buf[:]))

	if ShellScriptDir != "" {
		filename = filepath.Join(ShellScriptDir, filename)
	}

	err = ioutil.WriteFile(filename, data, 0700)
	return
}
