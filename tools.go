// +build tools

package main

import (
	// Import golangci-lint
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	// Import misspell
	_ "github.com/client9/misspell/cmd/misspell"
	// import goimports
	_ "golang.org/x/tools/cmd/goimports"
)
