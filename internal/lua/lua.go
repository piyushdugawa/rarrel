package lua

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"

	huh "charm.land/huh/v2"
)

func Prompt(L *lua.LState) int {
	lv := L.Get(-1)
	if str, ok := lv.(lua.LString); ok {
		// lv is LString
		var ret_val string
		huh.NewInput().
			Title(string(str)).
			Value(&ret_val).
			Run() // this is blocking...

		L.Push(lua.LString(ret_val))

	}
	if lv.Type() != lua.LTString {
		L.RaiseError("string required")
		return 0
	}
	return 1
}

func LoadLua(plugin string) {
	plugin_path := filepath.Join("./assets/plugins", plugin, "plugin.lua")
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("prompt", L.NewFunction(Prompt))
	if err := L.DoFile(plugin_path); err != nil {
		panic(err)
	}
}
