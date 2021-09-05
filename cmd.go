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

// Package exec executes a command in a new process.
package exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

var bufpool = sync.Pool{New: func() interface{} {
	return bytes.NewBuffer(make([]byte, 256))
}}

// DefaultTimeout is the default timeout.
var DefaultTimeout time.Duration

// CombinedOutputCmdHook is a CmdHook to run the command and returns its
// combined standard output and standard error, like exec.Cmd.CombinedOutput.
func CombinedOutputCmdHook(cmd *exec.Cmd) (stdout, stderr string, err error) {
	buf := bufpool.Get().(*bytes.Buffer)
	buf.Reset()

	cmd.Stdout = buf
	cmd.Stderr = buf
	err = cmd.Run()
	stdout = buf.String()
	stderr = buf.String()

	bufpool.Put(buf)
	return
}

// StdoutAndStderrBufferCmdHook returns a CmdHook that uses the given stdout
// and stderr buffer as the stdout and stderr of exec.Cmd.
func StdoutAndStderrBufferCmdHook(stdout, stderr *bytes.Buffer) CmdHook {
	return func(cmd *exec.Cmd) (_stdout, _stderr string, err error) {
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err = cmd.Run()
		_stdout = stdout.String()
		_stderr = stderr.String()
		return
	}
}

// Result represents the result of the executed command, which has implemented
// the interface error, so may be used as an error.
type Result struct {
	Name string
	Args []string

	Stdout string
	Stderr string
	Err    error
}

// NewResult returns a new Result.
func NewResult(name string, args []string, stdout, stderr string, err error) Result {
	return Result{Name: name, Args: args, Stdout: stdout, Stderr: stderr, Err: err}
}

// Error implements the interface error.
func (r Result) Error() string {
	err := r.Err.Error()
	buf := bufpool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString("cmd=")
	buf.WriteString(r.Name)
	if len(r.Args) > 0 {
		fmt.Fprintf(buf, ", args=%s", r.Args)
	}
	if len(r.Stdout) > 0 {
		buf.WriteString(", stdout=")
		buf.WriteString(r.Stdout)
	}
	if len(r.Stderr) > 0 {
		buf.WriteString(", stderr=")
		buf.WriteString(r.Stderr)
	}
	buf.WriteString(", err=")
	buf.WriteString(err)
	errmsg := buf.String()

	bufpool.Put(buf)
	return errmsg
}

// Unwrap implements errors.Unwrap.
func (r Result) Unwrap() error { return r.Err }

// CmdHook is used to customize how to run the command.
type CmdHook func(cmd *exec.Cmd) (stdout, stderr string, err error)

// Cmd represents a command executor.
type Cmd struct {
	// If not nil, it will be locked during the command is executed.
	//
	// Default: nil
	Lock *sync.Mutex

	// Shell is used to execute the command as the shell.
	//
	// If empty, use DefaultShell by default.
	Shell string

	// Timeout is used to produce the timeout context based on the context
	// argument if not 0 when executing the command.
	//
	// If empty, use DefaultTimeout by default.
	Timeout time.Duration

	// CmdHook is used to customize how to run the command.
	CmdHook CmdHook

	// ResultHook is used to to observe the result of the command.
	ResultHook func(Result)
}

// WithLock returns a new Cmd with the lock.
func (c Cmd) WithLock(lock *sync.Mutex) Cmd {
	c.Lock = lock
	return c
}

// WithShell returns a new Cmd with the shell.
func (c Cmd) WithShell(shell string) Cmd {
	c.Shell = shell
	return c
}

// WithTimeout returns a new Cmd with the timeout.
func (c Cmd) WithTimeout(timeout time.Duration) Cmd {
	c.Timeout = timeout
	return c
}

// WithCmdHook returns a new Cmd with the cmd hook.
func (c Cmd) WithCmdHook(hook CmdHook) Cmd {
	c.CmdHook = hook
	return c
}

// WithResultHook returns a new Cmd with the result hook.
func (c Cmd) WithResultHook(hook func(Result)) Cmd {
	c.ResultHook = hook
	return c
}

func (c Cmd) defaultCmdHook(cmd *exec.Cmd) (stdout, stderr string, err error) {
	stdoutbuf := bufpool.Get().(*bytes.Buffer)
	stderrbuf := bufpool.Get().(*bytes.Buffer)
	stdoutbuf.Reset()
	stderrbuf.Reset()

	cmd.Stdout = stdoutbuf
	cmd.Stderr = stderrbuf
	err = cmd.Run()
	stdout = stdoutbuf.String()
	stderr = stderrbuf.String()

	bufpool.Put(stdoutbuf)
	bufpool.Put(stderrbuf)
	return
}

func (c Cmd) runCmd(cmd *exec.Cmd) (stdout, stderr string, err error) {
	if c.Lock != nil {
		c.Lock.Lock()
		defer c.Lock.Unlock()
	}

	if c.CmdHook == nil {
		stdout, stderr, err = c.defaultCmdHook(cmd)
	} else {
		stdout, stderr, err = c.CmdHook(cmd)
	}

	return
}

// Run executes the command "name" with its arguments "args",
// then returns the stdout and stderr.
//
// Notice: if err is not nil, it is Result.
func (c Cmd) Run(cxt context.Context, name string, args ...string) (
	stdout, stderr string, err error) {
	if name == "" {
		panic("the cmd name is empty")
	}

	var cancel func()
	if c.Timeout > 0 {
		cxt, cancel = context.WithTimeout(cxt, c.Timeout)
		defer cancel()
	} else if DefaultTimeout > 0 {
		cxt, cancel = context.WithTimeout(cxt, DefaultTimeout)
		defer cancel()
	}

	cmd := exec.CommandContext(cxt, name, args...)
	stdout, stderr, err = c.runCmd(cmd)

	result := NewResult(name, args, stdout, stderr, err)
	if c.ResultHook != nil {
		c.ResultHook(result)
	}

	if err != nil {
		err = result
	}

	return
}

// Execute is the same as RunCmd, but only returns the error.
func (c Cmd) Execute(cxt context.Context, name string, args ...string) error {
	_, _, err := c.Run(cxt, name, args...)
	return err
}

// Output is the same as RunCmd, but only returns the stdout and the error.
func (c Cmd) Output(cxt context.Context, name string, args ...string) (string, error) {
	stdout, _, err := c.Run(cxt, name, args...)
	return stdout, err
}
