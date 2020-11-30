package pkg

import "fmt"

var (
	// BuildVersion of the package
	BuildVersion = "development"
	// BuildTime of the package
	BuildTime = "now"
	// UserAgent of the package
	UserAgent = fmt.Sprintf("gravl/%s (https://github.com/bzimmer/gravl)", BuildVersion)
	// PackageName is the name of the package
	PackageName = "github.com.bzimmer.gravl"
)
