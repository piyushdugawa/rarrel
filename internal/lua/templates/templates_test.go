package templates

import (
	"bytes"
	"testing"
	"text/template"

	lua "github.com/yuin/gopher-lua"
)

func TestToolsVersionCanBeUsedInTemplate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`
		Tools = {
			{
				name = "go",
				version = "1.24",
			},
		}
	`); err != nil {
		t.Fatal(err)
	}

	tmpl, err := template.New("go.mod").Parse(`go {{(index .Tools 0).version}}`)
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, collectGlobals(L)); err != nil {
		t.Fatal(err)
	}

	if got, want := output.String(), "go 1.24"; got != want {
		t.Fatalf("template output = %q, want %q", got, want)
	}
}
