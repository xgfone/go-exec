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

// RunShellScript is equal to DefaultCmd.RunShellScript(ctx, script, args...).
func RunShellScript(ctx context.Context, script string, args ...string) (string, string, error) {
	return DefaultCmd.RunShellScript(ctx, script, args...)
}

// RunShellCmd is equal to DefaultCmd.RunShellCmd(ctx, cmdfmt, cmdargs...).
func RunShellCmd(ctx context.Context, cmdfmt string, cmdargs ...string) (string, string, error) {
	return DefaultCmd.RunShellCmd(ctx, cmdfmt, cmdargs...)
}

// ExecuteShell is equal to DefaultCmd.ExecuteShell(ctx, cmdfmt, cmdargs...).
func ExecuteShell(ctx context.Context, cmdfmt string, cmdargs ...string) error {
	return DefaultCmd.ExecuteShell(ctx, cmdfmt, cmdargs...)
}

// OutputShell is equal to DefaultCmd.OutputShell(ctx, cmdfmt, cmdargs...).
func OutputShell(ctx context.Context, cmdfmt string, cmdargs ...string) (string, error) {
	return DefaultCmd.OutputShell(ctx, cmdfmt, cmdargs...)
}
