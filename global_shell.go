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

// ExecuteShellScript is equal to DefaultCmd.ExecuteShellScript.
func ExecuteShellScript(ctx context.Context, script string, args ...string) error {
	return DefaultCmd.ExecuteShellScript(ctx, script, args...)
}

// OutputShellScript is equal to DefaultCmd.OutputShellScript
func OutputShellScript(ctx context.Context, script string, args ...string) (
	stdout string, err error) {
	return DefaultCmd.OutputShellScript(ctx, script, args...)
}

// RunShellScript is equal to DefaultCmd.RunShellScript(ctx, script, args...).
func RunShellScript(ctx context.Context, script string, args ...string) (
	stdout, stderr string, err error) {
	return DefaultCmd.RunShellScript(ctx, script, args...)
}

// RunShellCmd is equal to DefaultCmd.RunShellCmd(ctx, cmdfmt, cmdargs...).
func RunShellCmd(ctx context.Context, cmdfmt string, cmdargs ...string) (
	stdout, stderr string, err error) {
	return DefaultCmd.RunShellCmd(ctx, cmdfmt, cmdargs...)
}

// ExecuteShellCmd is equal to DefaultCmd.ExecuteShellCmd(ctx, cmdfmt, cmdargs...).
func ExecuteShellCmd(ctx context.Context, cmdfmt string, cmdargs ...string) error {
	return DefaultCmd.ExecuteShellCmd(ctx, cmdfmt, cmdargs...)
}

// OutputShellCmd is equal to DefaultCmd.OutputShellCmd(ctx, cmdfmt, cmdargs...).
func OutputShellCmd(ctx context.Context, cmdfmt string, cmdargs ...string) (
	stdout string, err error) {
	return DefaultCmd.OutputShellCmd(ctx, cmdfmt, cmdargs...)
}
