package templates

import (
	"fmt"
	"path/filepath"

	"charm.land/huh/v2"
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

		entries := make([]struct {
			key   lua.LValue
			value lua.LValue
		}, 0, v.Len())
		isArray := true
		maxIndex := 0

		v.ForEach(func(key, value lua.LValue) {
			entries = append(entries, struct {
				key   lua.LValue
				value lua.LValue
			}{key: key, value: value})

			index, ok := key.(lua.LNumber)
			if !ok || index < 1 || index != lua.LNumber(int(index)) {
				isArray = false
				return
			}
			if int(index) > maxIndex {
				maxIndex = int(index)
			}
		})

		if len(entries) > 0 && isArray && maxIndex == len(entries) {
			result := make([]any, len(entries))
			for _, entry := range entries {
				index := int(entry.key.(lua.LNumber)) - 1
				result[index] = luaToGoSeen(entry.value, seen)
			}
			return result
		}

		allStringKeys := true
		for _, entry := range entries {
			if _, ok := entry.key.(lua.LString); !ok {
				allStringKeys = false
				break
			}
		}

		if allStringKeys {
			result := make(map[string]any, len(entries))
			for _, entry := range entries {
				result[string(entry.key.(lua.LString))] = luaToGoSeen(entry.value, seen)
			}
			return result
		}

		// Preserve mixed tables by collecting all keys.
		result := make(map[any]any, len(entries))
		for _, entry := range entries {
			result[luaToGoSeen(entry.key, seen)] = luaToGoSeen(entry.value, seen)
		}
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

	var confirm bool

	huh.NewConfirm().
		Title("Are you sure?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&confirm).
		Run()

	if confirm {
		if err := generate.RenderTemplates(srcDir, destDir, vars); err != nil {
			panic(err)
		}
		return 0
	}

	fmt.Println("Operation cancelled by User!!!")

	return 1
}
