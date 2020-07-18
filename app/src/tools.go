// +build tools
// this package is not used, for ignoreing `go get -u some-go-tools` will be removed by go mod tidy...

package main

import (
	_ "github.com/uber/go-torch"
	_ "golang.org/x/tools/cmd/goimports"
)
