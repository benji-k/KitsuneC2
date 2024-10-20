//This is a "hacky" package. Because of the way Go's "embed" package works, all files that need to be embedded into a final compiled binary, need
//to reside in the same local directory that the go:embed directives are in.

//This package is only used by server/builder/builder.go to get access to the implants source code.

package implantSource

import (
	"embed"
)

var (
	//go:embed implant/**
	ImplantFs embed.FS

	//go:embed lib/**
	LibFs embed.FS

	//go:embed go.sum
	GoSum string

	//go:embed go.mod
	GoMod string
)
