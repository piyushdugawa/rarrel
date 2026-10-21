package lua

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func Double(L *lua.LState) int {
	lv := L.ToInt(1)            /* get argument */
	L.Push(lua.LNumber(lv * 2)) /* push result */
	return 1                    /* number of results */
}

func LoadLua(plugin string) {
	plugin_path := filepath.Join("./assets/plugins", plugin, "plugin.lua")
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("double", L.NewFunction(Double)) /* Original lua_setglobal uses stack... */
	if err := L.DoFile(plugin_path); err != nil {
		panic(err)
	}
}
