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

package exec

import (
	"context"
	"sync"
	"time"
)

// DefaultCmd is the global default cmd executor.
var DefaultCmd Cmd

func init() {
	DefaultCmd.Lock = new(sync.Mutex)
	DefaultCmd.ResultHook = LogResult
}

// WithLock is equal to DefaultCmd.WithLock(lock).
func WithLock(lock *sync.Mutex) Cmd {
	return DefaultCmd.WithLock(lock)
}

// WithShell is equal to DefaultCmd.WithShell(shell).
func WithShell(shell string) Cmd {
	return DefaultCmd.WithShell(shell)
}

// WithTimeout is equal to DefaultCmd.WithTimeout(timeout).
func WithTimeout(timeout time.Duration) Cmd {
	return DefaultCmd.WithTimeout(timeout)
}

// WithCmdHook is equal to DefaultCmd.WithCmdHook(hook).
func WithCmdHook(hook CmdHook) Cmd {
	return DefaultCmd.WithCmdHook(hook)
}

// WithResultHook is equal to DefaultCmd.WithResultHook(hook).
func WithResultHook(hook func(Result)) Cmd {
	return DefaultCmd.WithResultHook(hook)
}

// Run is equal to DefaultCmd.Run(ctx, name, args...).
func Run(ctx context.Context, name string, args ...string) (stdout, stderr string, err error) {
	return DefaultCmd.Run(ctx, name, args...)
}

// Execute is equal to DefaultCmd.Execute(cxt, name, args...).
func Execute(cxt context.Context, name string, args ...string) error {
	return DefaultCmd.Execute(cxt, name, args...)
}

// Output is equal to DefaultCmd.Output(cxt, name, args...).
func Output(cxt context.Context, name string, args ...string) (string, error) {
	return DefaultCmd.Output(cxt, name, args...)
}
