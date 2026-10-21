package lua

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func LoadLua(plugin string) {
	plugin_path := filepath.Join("./assets/plugins", plugin, "plugin.lua")
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(plugin_path); err != nil {
		panic(err)
	}
}
