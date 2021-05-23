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

import "context"

// DefaultCmd is the global default cmd executor.
var DefaultCmd = NewCmd()

// AppendResultHooks is equal to DefaultCmd.AppendResultHooks(hooks...).
func AppendResultHooks(hooks ...ResultHook) {
	DefaultCmd.AppendResultHooks(hooks...)
}

// RunCmd is equal to DefaultCmd.RunCmd(ctx, name, args...).
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
