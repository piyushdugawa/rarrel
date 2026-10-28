package templates

import (
	"path/filepath"

	"github.com/piyushdugawa/rarrel/internal/config"
	"github.com/piyushdugawa/rarrel/internal/generate"
	lua "github.com/yuin/gopher-lua"
)

func luaToGo(v lua.LValue) any {
	return luaToGoSeen(v, make(map[*lua.LTable]bool))
}

func luaToGoSeen(v lua.LValue, seen map[*lua.LTable]bool) any {
	switch v := v.(type) {
	case lua.LString:
		return string(v)

	case lua.LNumber:
		return float64(v)

	case lua.LBool:
		return bool(v)

	case *lua.LTable:
		if seen[v] {
			return "<circular reference>"
		}

		seen[v] = true
		defer delete(seen, v)

		// Preserve mixed tables by collecting all keys.
		result := make(map[any]any)

		v.ForEach(func(key, value lua.LValue) {
			result[luaToGoSeen(key, seen)] = luaToGoSeen(value, seen)
		})

		return result

	case *lua.LFunction:
		return nil

	default:
		return nil
	}
}

func collectGlobals(L *lua.LState) map[string]any {
	vars := make(map[string]any)

	L.G.Global.ForEach(func(key, value lua.LValue) {
		name, ok := key.(lua.LString)
		if !ok {
			return
		}

		// Exclude functions and built-in globals.
		if value.Type() == lua.LTFunction {
			return
		}

		vars[string(name)] = luaToGo(value)
	})

	return vars
}

func RenderTemplate(L *lua.LState) int {
	LsrcDir := L.Get(1)
	srcDir := filepath.Join(config.PluginDir, LsrcDir.String())
	LdestDir := L.Get(2)
	destDir := LdestDir.String()

	vars := collectGlobals(L)

	if err := generate.RenderTemplates(srcDir, destDir, vars); err != nil {
		panic(err)
	}

	return 0
}
