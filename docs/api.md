# Lua API Reference

This document provides a detailed specification of all Go functions and runtime variables exposed to the Lua scripting environment in **Rarrel** (`rrl`).

---

## Table of Contents

- [Overview](#overview)
- [Runtime Environment & Globals](#runtime-environment--globals)
  - [`PLUGIN_DIR`](#plugin_dir)
  - [`package.path`](#packagepath)
- [Exposed Go Functions](#exposed-go-functions)
  - [`prompt(title, prefix, placeholder)`](#prompttitle-prefix-placeholder)
  - [`verifyTools(tools)`](#verifytoolstools)
  - [`renderTemplates(srcDir, destDir)`](#rendertemplatessrcdir-destdir)
- [Type Marshalling (Lua to Go)](#type-marshalling-lua-to-go)
- [Best Practices for Plugin Scripts](#best-practices-for-plugin-scripts)

---

## Overview

Rarrel embeds the [GopherLua](https://github.com/yuin/gopher-lua) runtime inside the Go core. When a plugin is executed via `rrl create <plugin>`, Rarrel opens standard Lua libraries (`io`, `os`, `string`, `table`, etc.) and injects custom Go functions into the Lua global scope.

```
+--------------------------------------------------------------+
|                         Lua Plugin                           |
|       prompt()        verifyTools()       renderTemplates()  |
+-----------+-----------------+---------------------+----------+
            |                 |                     |
+-----------v-----------------v---------------------v----------+
|                       Go Core Engine                         |
|  charm.land/huh/v2     exec.LookPath      generate.Render()  |
|  (Terminal UI Input)   (System Binary)    (Go text/template) |
+--------------------------------------------------------------+
```

---

## Runtime Environment & Globals

Before loading `plugin.lua`, Rarrel automatically configures the following global variables:

### `PLUGIN_DIR`

- **Type**: `string`
- **Description**: The absolute directory path of the plugin currently being executed.
- **Example**: `"/home/user/.rarrel/plugins/go_starter"`
- **Usage**: Useful for resolving local assets or relative files within the plugin folder.

### `package.path`

Rarrel automatically prepends the plugin directory to Lua's module search path:

```lua
package.path = PLUGIN_DIR .. "/?.lua;" ..
               PLUGIN_DIR .. "/?/plugin.lua;" ..
               package.path
```

This enables plugins to be modular and split across multiple Lua files:

```lua
-- Loads ~/.rarrel/plugins/<plugin>/helpers.lua
local helpers = require("helpers")
```

---

## Exposed Go Functions

### `prompt(title, prefix, placeholder)`

Displays a styled interactive terminal input prompt to the user and returns the entered string.

- **Go Implementation**: `internal/lua.Prompt`
- **Underlying Library**: [Charm Huh (`huh.NewInput()`)](https://github.com/charmbracelet/huh)

#### Signature

```lua
local input = prompt(title, prefix, placeholder)
```

#### Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `title` | `string` | **Yes** | The question or prompt heading displayed above the input. |
| `prefix` | `string` | **Yes** | The prompt indicator displayed directly before the input field (e.g., `"? "` or `"> "`). |
| `placeholder` | `string` | **Yes** | Default fallback text. If the user presses `Enter` without typing, this value is returned. |

#### Returns

- `string`: The trimmed text entered by the user, or `placeholder` if the user provided an empty string.

#### Error Handling

- If `title` is not of type `string`, a Lua error is raised (`string required`).

#### Example

```lua
-- Prompt user for project name, defaulting to current folder name
local pwd = io.popen("pwd"):read("*l")
local basename = pwd:match("([^/]+)/*$")

ProjectName = prompt("Enter Project Name:", "? ", basename)
ModulePath  = prompt("Enter Module Path:", "? ", "github.com/user/" .. ProjectName)
```

---

### `verifyTools(tools)`

Checks whether specified command-line tools are installed in the host system's `$PATH`, executes each tool to determine its installed version, and compares it against minimum required versions.

- **Go Implementation**: `internal/lua/tools.VerifyTools`

#### Signature

```lua
verifyTools(toolsTable)
```

#### Parameters

- **`toolsTable`** (`table` of tables): An array of tool definition tables with the following fields:

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `name` | `string` | **Yes** | The binary name to check in `$PATH` (e.g., `"go"`, `"git"`, `"docker"`). |
| `required` | `boolean` | **Yes** | If `true`, a missing tool or failing version constraint terminates execution immediately. |
| `version` | `string` | No | Minimum version required (e.g., `"1.24"`). If omitted (`nil`), any version is accepted. |
| `version_arg` | `string` | **Yes** | Command-line argument to run to fetch version info (e.g., `"version"` or `"--version"`). |

#### Returns

- None (`0` values returned to Lua).

#### How It Works Internally

1. **Existence Check**: Uses Go's `os/exec.LookPath(tool.Name)` to verify binary existence.
2. **Version Extraction**: Executes `exec.Command(tool.Name, tool.VerArgs).CombinedOutput()`.
3. **Regex Matching**: Applies regex pattern `\d+(?:\.\d+){1,2}` against output to extract semantic versions.
4. **Validation Logic**:
   - **Inherited Version**: If `version` is not defined (`nil`), any installed version passes and prints:
     ```text
     git v2.43.0 (Inherited) ✓
     ```
   - **Version Constraint**: If `version` is defined, checks if `extractedVersion >= tool.Version`:
     - If passing:
       ```text
       go v1.24.0 (>=1.24) ✓
       ```
     - If failing:
       ```text
       go v1.20.0 found but go v1.24 is required ✘
       ```
       If `required = true`, the program halts.
   - **Missing Tools**:
     - If `required = false`:
       ```text
       jq ✘  --Ignored
       ```
     - If `required = true`:
       ```text
       docker ✘
       ```
       Halts execution immediately.

#### Example

```lua
Tools = {
    {
        name = "go",
        required = true,
        version = "1.24",
        version_arg = "version"
    },
    {
        name = "git",
        required = true,
        version_arg = "--version"
    },
    {
        name = "docker",
        required = false,
        version_arg = "--version"
    }
}

verifyTools(Tools)
```

---

### `renderTemplates(srcDir, destDir)`

Marshals all global variables in Lua into a Go template context, requests interactive user confirmation, and recursively renders `.tmpl` files to the target destination.

- **Go Implementation**: `internal/lua/templates.RenderTemplate`
- **Template Engine**: Go standard library `text/template`

#### Signature

```lua
local status = renderTemplates(srcDir, destDir)
```

#### Parameters

| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `srcDir` | `string` | **Yes** | Relative path to the template folder from the plugin directory (e.g., `"./templates"`). |
| `destDir` | `string` | **Yes** | Destination path on disk where files should be rendered (e.g., `"./"` or `"./my-project"`). |

#### Returns

- `number`: `0` on successful generation, `1` if cancelled by user.

#### Internal Execution Flow

1. **State Gathering (`collectGlobals`)**:
   Traverses all variables in Lua's global table (`L.G.Global`). Filters out built-in functions, and recursively maps Lua types to native Go types.
2. **Interactive Confirmation**:
   Prompts the user via `huh.NewConfirm()`:
   ```text
   Are you sure?
   > Yes! / No.
   ```
   If the user selects "No.", generation aborts with:
   ```text
   Operation cancelled by User!!!
   ```
3. **Template Walking & Rendering**:
   - Calls `generate.RenderTemplates(srcDir, destDir, vars)`.
   - Uses `filepath.WalkDir` to iterate over files in `srcDir`.
   - Ignores non-`.tmpl` files.
   - Preserves directory structures and creates parent directories via `os.MkdirAll(..., 0755)`.
   - Trims the `.tmpl` extension to determine the final output filename (`main.go.tmpl` $\rightarrow$ `main.go`).
   - Parses the template and executes it with the collected globals passed as root `.`.

#### Example

```lua
-- All global variables are passed directly to template rendering!
ProjectName = "my-service"
ModulePath  = "github.com/user/my-service"
Port        = 8080
EnableDocs  = true

renderTemplates("./templates", "./" .. ProjectName)
```

---

## Type Marshalling (Lua to Go)

When `renderTemplates` gathers globals, the Go core maps Lua values to native Go data structures as follows:

| Lua Type | Marshalled Go Type | Example Lua Value | Template Access |
| :--- | :--- | :--- | :--- |
| `string` | `string` | `Name = "app"` | `{{.Name}}` |
| `number` | `float64` | `Port = 8080` | `{{.Port}}` |
| `boolean` | `bool` | `Debug = true` | `{{if .Debug}}...{{end}}` |
| Sequential Table (Array) | `[]any` | `Tags = {"api", "v1"}` | `{{index .Tags 0}}`, `{{range .Tags}}` |
| Associative Table (Map) | `map[string]any` | `Meta = {env = "prod"}` | `{{.Meta.env}}` |
| Nested Table | `[]any` or `map[string]any` | `Tools = {{name = "go"}}` | `{{(index .Tools 0).name}}` |

*Note: Functions (`lua.LTFunction`) and circular references are automatically skipped or sanitized.*

---

## Best Practices for Plugin Scripts

1. **Use Globals for Template Data**: Only define global variables (`UpperCamelCase` or `PascalCase`) for data that templates will read. Use `local` variables for internal script logic.
2. **Always Run Pre-Flight Checks**: Call `verifyTools()` before `renderTemplates()` so dependencies fail fast before prompts or files are written.
3. **Handle Destination Gracefully**: Support in-place generation (`./`) when the user wants to scaffold into the current folder.
