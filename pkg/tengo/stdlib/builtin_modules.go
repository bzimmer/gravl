package stdlib

import (
	"github.com/d5/tengo/v2"
)

// BuiltinModules are builtin type standard library modules
var BuiltinModules = map[string]map[string]tengo.Object{
	"strava": stravaModule,
	"gravl":  gravlModule,
}
