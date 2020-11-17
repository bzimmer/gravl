package stdlib

import (
	"github.com/d5/tengo/v2"
)

// AllModuleNames returns a list of all default module names
func AllModuleNames() []string {
	var names []string
	for name := range BuiltinModules {
		names = append(names, name)
	}
	return names
}

// GetModuleMap returns the module map of all modules for the given module names
func GetModuleMap(names ...string) *tengo.ModuleMap {
	modules := tengo.NewModuleMap()
	for _, name := range names {
		if mod := BuiltinModules[name]; mod != nil {
			modules.AddBuiltinModule(name, mod)
		}
	}
	return modules
}
