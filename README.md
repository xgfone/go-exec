# Command Execution [![Build Status](https://github.com/xgfone/go-exec/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-exec/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-exec)](https://pkg.go.dev/github.com/xgfone/go-exec) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-exec/master/LICENSE)

The package supplies the common command execution supporting `Go1.7+`.

## Install
```shell
$ go get -u github.com/xgfone/go-exec
```

## Example
```go
package main

import (
	"context"
	"fmt"

	"github.com/xgfone/go-exec"
)

const scripttmpl = `
ls %s
rm -rf %s
`

func main() {
	exec.Execute(context.Background(), "mkdir", "testdir")
	exec.ExecuteShellCmd(context.Background(), "echo abc > %s/%s", "testdir", "testfile")

	data, _ := exec.OutputShellCmd(context.Background(), "cat %s/%s", "testdir", "testfile")
	fmt.Println(data)

	_, _, err := exec.RunShellScript(context.Background(), scripttmpl, "testdir", "testdir")
	fmt.Println(err)

	// Output:
	// abc
	//
	// <nil>
}
```
