package tools

import (
	"fmt"
	"os/exec"

	lua "github.com/yuin/gopher-lua"
)

func ToolExist(L *lua.LState) int {
	Ltool := L.Get(1)
	tool := Ltool.String()
	path, err := exec.LookPath(tool)
	if err != nil {
		fmt.Println("Go is not installed or not in PATH!")
		return 0
	}
	fmt.Printf("%s executable: %s\n", tool, path)
	return 0
}
