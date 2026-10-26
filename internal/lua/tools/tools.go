package tools

import (
	"fmt"
	"os/exec"
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

type Tool struct {
	Name     string
	Version  *string
	VerArgs  string
	Required bool
}

var tools []Tool

func VerifyTools(L *lua.LState) int {
	lv := L.Get(1)
	if tbl, ok := lv.(*lua.LTable); ok {
		for i := 1; i <= L.ObjLen(tbl); i++ {
			toolTbl := tbl.RawGetInt(i).(*lua.LTable)

			var tool Tool
			tool.Name = toolTbl.RawGetString("name").String()

			if v := toolTbl.RawGetString("version"); v != lua.LNil {
				if s, ok := v.(lua.LString); ok {
					version := string(s)
					tool.Version = &version
				}
			}

			tool.VerArgs = toolTbl.RawGetString("version_arg").String()
			b, err := LuaValueToBool(toolTbl.RawGetString("required"))
			if err != nil {
				fmt.Println("Error:", err)
			}
			tool.Required = b
			tools = append(tools, tool)
		}
		process(tools)
	}
	return 0
}

func process(tools []Tool) {
	fmt.Println()
	for _, tool := range tools {

		_, err := exists(tool.Name)
		if err == nil {
			output, err := exec.Command(tool.Name, tool.VerArgs).CombinedOutput()
			if err != nil {
				fmt.Println("Error:", err)
			}

			extractedVersion := extractVersion(string(output))

			if tool.Version == nil {
				fmt.Printf("%s \033[1mv%s\033[0m(\033[1;34mInherited\033[0m) \033[1;34m✓\033[0m\n", tool.Name, extractedVersion)
			} else {
				if extractedVersion > *tool.Version || extractedVersion == *tool.Version {
					fmt.Printf("%s \033[1mv%s\033[0m (\033[1m>=%s\033[0m) \033[1;32m✓\033[0m\n", tool.Name, extractedVersion, *tool.Version)
				} else {
					fmt.Printf("%s \033[1mv%s\033[0m found but %s \033[1mv%s\033[0m is required \033[1;31m\u2718\033[0m\n", tool.Name, extractedVersion, tool.Name, *tool.Version)
					return
				}
			}
		}

		if err != nil {
			if !tool.Required {
				fmt.Printf("%s \033[1;31m\u2718\033[0m  \033[1;33m--Ignored\033[0m\n", tool.Name)
			} else {
				fmt.Printf("%s \033[1;31m\u2718\033[0m\n", tool.Name)
				return
			}
		}

	}
}

func LuaValueToBool(v lua.LValue) (bool, error) {
	b, ok := v.(lua.LBool)
	if !ok {
		return false, fmt.Errorf("expected boolean, got %s", v.Type())
	}
	return bool(b), nil
}

func exists(tool string) (string, error) {
	path, err := exec.LookPath(tool)
	if err != nil {
		return "", err
	}
	return path, nil
}

var versionRegex = regexp.MustCompile(`\d+(?:\.\d+){1,2}`)

func extractVersion(output string) string {
	return versionRegex.FindString(output)
}
