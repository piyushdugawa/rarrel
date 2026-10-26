package tools

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Tool struct {
	Name     string
	Version  string
	VerArgs  string
	Required bool
}

// var tools []Tool

func VerifyTools(L *lua.LState) int {
	lv := L.Get(1) // get the value at the top of the stack
	if tbl, ok := lv.(*lua.LTable); ok {
		// lv is LTable
		fmt.Println(L.ObjLen(tbl))
		for i := 1; i <= L.ObjLen(tbl); i++ {
			toolTbl := tbl.RawGetInt(i).(*lua.LTable)

			var tool Tool
			tool.Name = toolTbl.RawGetString("name").String()
			tool.Version = toolTbl.RawGetString("version").String()
			tool.VerArgs = toolTbl.RawGetString("version_arg").String()
			b, err := LuaValueToBool(toolTbl.RawGetString("required"))
			if err != nil {
				// handle error
			}
			tool.Required = b
			// tools = append(tools, tool)
			exists(tool.Name)

			output, err := exec.Command(tool.Name, tool.VerArgs).CombinedOutput()
			if err != nil {
				L.RaiseError("%s", err.Error())
				return 0
			}

			fmt.Println(extractVersion(string(output)))

		}
	}
	return 0
}

func LuaValueToBool(v lua.LValue) (bool, error) {
	b, ok := v.(lua.LBool)
	if !ok {
		return false, fmt.Errorf("expected boolean, got %s", v.Type())
	}
	return bool(b), nil
}

// func ToolExist(L *lua.LState) int {
// 	Ltool := L.Get(1)
// 	tool := Ltool.String()

// 	path, err := exists(tool)
// 	if err != nil {
// 		L.Push(lua.LFalse)
// 		return 1
// 	}

// 	L.Push(lua.LString(path))
// 	L.Push(lua.LTrue)
// 	return 2

// }

func exists(tool string) (string, error) {
	path, err := exec.LookPath(tool)
	if err != nil {
		fmt.Printf("%s is either not installed or not in PATH\n", tool)
		return "", err
	}
	fmt.Printf("%s found at %s\n", tool, path)
	return path, nil
}

func ChkInstalledVer(L *lua.LState) int {
	LCmd := L.Get(1)
	cmnd := strings.Fields(LCmd.String())
	output, err := exec.Command(cmnd[0], cmnd[1.]).CombinedOutput()
	if err != nil {
		L.RaiseError("%s", err.Error())
		return 0
	}

	L.Push(lua.LString(extractVersion(string(output))))
	return 1
}

var versionRegex = regexp.MustCompile(`\d+(?:\.\d+){1,2}`)

func extractVersion(output string) string {
	return versionRegex.FindString(output)
}
