package lua

import (
	"path/filepath"

	"github.com/piyushdugawa/rarrel/internal/lua/tools"
	lua "github.com/yuin/gopher-lua"

	huh "charm.land/huh/v2"
)

func Prompt(L *lua.LState) int {
	title := L.Get(1)
	prompt := L.Get(2)
	placeholder := L.Get(3)
	if title_str, ok := title.(lua.LString); ok {
		var ret_val string
		huh.NewInput().
			Title(string(title_str)).
			Prompt(string(prompt.(lua.LString))).
			Placeholder(string(placeholder.(lua.LString))).
			Value(&ret_val).
			Run()

		L.Push(lua.LString(ret_val))

	}
	if title.Type() != lua.LTString {
		L.RaiseError("string required")
		return 0
	}
	return 1
}

func LoadLua(plugin string) {
	plugin_path := filepath.Join("./assets/plugins", plugin, "plugin.lua")
	plugin_dir := filepath.Dir(plugin_path)

	L := lua.NewState()
	defer L.Close()

	L.OpenLibs()

	L.SetGlobal("PLUGIN_DIR", lua.LString(plugin_dir))
	L.SetGlobal("prompt", L.NewFunction(Prompt))
	L.SetGlobal("verifyTools", L.NewFunction(tools.VerifyTools))

	if err := L.DoString(`
	package.path = PLUGIN_DIR .. "/?.lua;" ..
	PLUGIN_DIR .. "/?/plugin.lua;" ..
	package.path`); err != nil {
		panic(err)
	}

	if err := L.DoFile(plugin_path); err != nil {
		panic(err)
	}
}
