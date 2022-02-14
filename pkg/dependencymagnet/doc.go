//go:build tools
// +build tools

package dependencymagnet

// Include dependencies that would otherwise not get imported
import (
	// Used to generate bindata
	_ "github.com/jteeuwen/go-bindata/go-bindata"
)
