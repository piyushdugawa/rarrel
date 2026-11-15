package ui

import (
	"fmt"

	"charm.land/huh/v2"
	lua "github.com/yuin/gopher-lua"
)

func Selection(title string, options *lua.LTable, val *string) {

	// slice of options
	var opts []huh.Option[string]

	// Iterate over the outer table
	options.ForEach(func(key, value lua.LValue) {
		innerTable, ok := value.(*lua.LTable)
		if !ok {
			fmt.Println("Some error")
			return
		}

		// Read the inner table values
		k := innerTable.RawGetInt(1).String()
		v := innerTable.RawGetInt(2).String()

		opts = append(opts, huh.NewOption(k, v))
	})

	// Selection Menu
	huh.NewSelect[string]().
		Title(title).
		Options(
			opts...,
		).
		Value(val).
		Run()
}
