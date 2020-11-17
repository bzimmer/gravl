package stdlib

import (
	"github.com/d5/tengo/v2"

	"github.com/bzimmer/gravl/pkg"
)

var (
	gravlModule = map[string]tengo.Object{
		"version": version,
	}
)

var version = &tengo.UserFunction{
	Name: "version",
	Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 0 {
			err = tengo.ErrWrongNumArguments
			return
		}
		ret, err = tengo.FromInterface(pkg.BuildVersion)
		return
	},
}
