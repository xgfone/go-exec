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

// CmdError represents a cmd error.
type CmdError struct {
	Name string
	Args []string

	Err    error
	Stdout string
	Stderr string
}

// NewCmdError returns a new CmdError.
func NewCmdError(name string, args []string, stdout, stderr string, err error) CmdError {
	return CmdError{Name: name, Args: args, Stdout: stdout, Stderr: stderr, Err: err}
}

func (c CmdError) Error() string {
	err := c.Err.Error()
	buf := bytes.NewBuffer(nil)
	buf.Grow(len(c.Stderr) + len(c.Stdout) + len(err) + len(c.Name) + 36)
	buf.WriteString("cmd=")
	buf.WriteString(c.Name)

	if len(c.Args) > 0 {
		fmt.Fprintf(buf, ", args=%s", c.Args)
	}
	if len(c.Stdout) > 0 {
		buf.WriteString(", stdout=")
		buf.WriteString(c.Stdout)
	}
	if len(c.Stderr) > 0 {
		buf.WriteString(", stderr=")
		buf.WriteString(c.Stderr)
	}

	buf.WriteString(", err=")
	buf.WriteString(err)
	return buf.String()
}

// Unwrap implements errors.Unwrap.
func (c CmdError) Unwrap() error {
	return c.Err
}

// ResultHook is a hook to observe the result of the command.
type ResultHook func(name string, args []string, stdout, stderr string, err error)

// Cmd represents a command executor.
type Cmd struct {
	// ResultHooks is the hooks to observe the result of the command.
	ResultHooks []ResultHook

	// Shell is used to execute the command as the shell.
	//
	// If empty, use DefaultShell by default.
	Shell string

	// Timeout is used to produce the timeout context based on the context
	// argument if not 0 when executing the command.
	//
	// If empty, use DefaultTimeout by default.
	Timeout time.Duration

	// RunCmd allows the user to decide how to run the command.
	//
	// Default: cmd.Run()
	RunCmd func(cmd *exec.Cmd) error

	// If not nil, it will be locked during the command is executed.
	//
	// Default: nil
	Lock *sync.Mutex
}

// NewCmd returns a new executor Cmd.
func NewCmd() *Cmd { return new(Cmd) }

// AppendResultHooks appends some result hooks.
func (c *Cmd) AppendResultHooks(hooks ...ResultHook) {
	c.ResultHooks = append(c.ResultHooks, hooks...)
}

func (c *Cmd) runCmd(cmd *exec.Cmd) error {
	if c.Lock != nil {
		c.Lock.Lock()
		defer c.Lock.Unlock()
	}

	if c.RunCmd != nil {
		return c.RunCmd(cmd)
	}
	return cmd.Run()
}

func (c *Cmd) runResultHooks(name string, args []string, stdout, stderr string,
	err error) (string, string, error) {
	for _, hook := range c.ResultHooks {
		hook(name, args, stdout, stderr, err)
	}

	switch err.(type) {
	case nil, CmdError:
	default:
		err = NewCmdError(name, args, stdout, stderr, err)
	}

	return stdout, stderr, err
}

// Run executes the command, name, with its arguments, args,
// then returns stdout, stderr and error.
//
// Notice: if there is an error to be returned, it is CmdError.
func (c *Cmd) Run(cxt context.Context, name string, args ...string) (
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

	output := bufpool.Get().(*bytes.Buffer)
	errput := bufpool.Get().(*bytes.Buffer)
	output.Reset()
	errput.Reset()

	cmd := exec.CommandContext(cxt, name, args...)
	cmd.Stdout = output
	cmd.Stderr = errput
	err = c.runCmd(cmd)
	stdout = output.String()
	stderr = errput.String()

	bufpool.Put(output)
	bufpool.Put(errput)

	return c.runResultHooks(name, args, stdout, stderr, err)
}

// Execute is the same as RunCmd, but only returns the error.
func (c *Cmd) Execute(cxt context.Context, name string, args ...string) error {
	_, _, err := c.Run(cxt, name, args...)
	return err
}

// Output is the same as RunCmd, but only returns the stdout and the error.
func (c *Cmd) Output(cxt context.Context, name string, args ...string) (string, error) {
	stdout, _, err := c.Run(cxt, name, args...)
	return string(stdout), err
}
