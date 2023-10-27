// Copyright 2023 xgfone
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

import "log"

var logResult = _logresult

// LogResult is a hook to log the command result.
func LogResult(r Result) { logResult(r) }

func _logresult(r Result) {
	if r.Err == nil {
		log.Printf("successfully execute the command: cmd=%s, args=%v", r.Name, r.Args)
	} else {
		log.Printf("fail to execute the command: cmd=%s, args=%v, stdout=%s, stderr=%s, err=%s",
			r.Name, r.Args, r.Stdout, r.Stderr, r.Err.Error())
	}
}
