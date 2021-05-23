# Command Execution [![Build Status](https://travis-ci.org/xgfone/go-exec.svg?branch=master)](https://travis-ci.org/xgfone/go-exec) [![GoDoc](https://godoc.org/github.com/xgfone/go-exec?status.svg)](http://pkg.go.dev/github.com/xgfone/go-exec) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-exec/master/LICENSE)

The package supplies the common command execution.

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
	exec.ExecuteShell(context.Background(), "echo abc > %s/%s", "testdir", "testfile")

	data, _ := exec.OutputShell(context.Background(), "cat %s/%s", "testdir", "testfile")
	fmt.Println(data)

	_, _, err := exec.RunShellScript(context.Background(), scripttmpl, "testdir", "testdir")
	fmt.Println(err)

	// Output:
	// abc
	//
	// <nil>
}
```
