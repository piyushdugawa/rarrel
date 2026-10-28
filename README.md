# Rarrel (`rrl`)

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Lua Support](https://img.shields.io/badge/Lua-5.1%20%2F%20GopherLua-000080?style=flat&logo=lua)](https://github.com/yuin/gopher-lua)
[![CLI](https://img.shields.io/badge/CLI-Cobra-blue?style=flat)](https://github.com/spf13/cobra)
[![Terminal UI](https://img.shields.io/badge/TUI-Charm%20Huh-orange?style=flat)](https://github.com/charmbracelet/huh)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**Rarrel** (`rrl`) is a lightweight, extensible CLI project generator. It pairs the dynamic scripting power of **Lua plugins** with the speed, robust file I/O, interactive terminal prompts, and template rendering of **Go**.

---

## Highlights

- 🧩 **Lua-Powered Plugins**: Define project metadata, input prompts, and generation logic without recompiling the CLI.
- 💬 **Interactive Prompts**: Clean terminal forms with placeholders and defaults powered by Charm's [Huh](https://github.com/charmbracelet/huh).
- 🛠️ **Pre-Flight Tool Checks**: Validate host binaries (`go`, `git`, `docker`, etc.) and minimum versions before writing files.
- 📂 **Recursive `.tmpl` Rendering**: Automatically walks template folders, strips `.tmpl` extensions, creates directories, and renders using Go `text/template`.
- 🛡️ **Interactive Confirmation**: Prompts for confirmation before generating to prevent accidental file overwrites.

---

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/piyushdugawa/rarrel.git
cd rarrel

# Build and install to ~/go/bin
just install
# or manually: go build -o build/rrl ./cmd/rrl/
```

### Usage

```bash
# Generate a project using an installed plugin
rrl create <plugin-name>

# Check installed version
rrl version
```

---

## How It Works

Plugins reside in `~/.rarrel/plugins/<plugin_name>/` and contain a `plugin.lua` script plus a `templates/` folder:

```
~/.rarrel/plugins/go_starter/
├── plugin.lua
└── templates/
    ├── go.mod.tmpl
    └── main.go.tmpl
```

### 1. `plugin.lua`

```lua
-- Interactive prompt (backed by Charm Huh)
ProjectName = prompt("Enter Project Name:", "? ", "my-app")
ModulePath  = prompt("Enter Module Path:", "? ", "github.com/user/" .. ProjectName)

-- Validate dependencies
verifyTools({
    { name = "go", required = true, version = "1.24", version_arg = "version" },
    { name = "git", required = true, version_arg = "--version" }
})

-- Render templates with collected Lua globals
renderTemplates("./templates", "./" .. ProjectName)
```

### 2. `templates/go.mod.tmpl`

```gotemplate
module {{.ModulePath}}

go {{(index .Tools 0).version}}
```

Running `rrl create go_starter` prompts for inputs, verifies your tools, confirms with you, and outputs a ready-to-use project.

---

## Lua API at a Glance

Rarrel exposes the following Go functions to Lua:

| Function | Description |
| :--- | :--- |
| [`prompt(title, prefix, placeholder)`](docs/api.md#prompttitle-prefix-placeholder) | Displays an interactive terminal input prompt with default fallback. |
| [`verifyTools(toolsTable)`](docs/api.md#verifytoolstools) | Checks for CLI tools in `$PATH` and verifies version requirements. |
| [`renderTemplates(srcDir, destDir)`](docs/api.md#rendertemplatessrcdir-destdir) | Collects Lua globals, asks for user confirmation, and renders `.tmpl` files. |

👉 **Read the full [Lua API Reference](docs/api.md) for detailed signatures, parameter tables, and type conversions.**

---

## Documentation

Comprehensive guides and deep dives are available in the [`docs/`](docs/) directory:

- 📖 **[Lua API Reference](docs/api.md)**: Deep dive into all Go functions exposed to Lua, type marshalling, and runtime variables (`PLUGIN_DIR`, `package.path`).
- 🛠️ **[Plugin & Template Guide](docs/plugins.md)**: Tutorial on creating custom plugins, directory layout, and Go `text/template` syntax.
- 🏗️ **[Architecture & Internals](docs/architecture.md)**: How Rarrel bridges Go and Lua, handles the TUI lifecycle, and manages recursive file generation.

---

## Development

```bash
just fmt      # Format Go code
just lint     # Run golangci-lint
just test     # Run unit tests
just build    # Build binary
```

---

## License

This project is licensed under the [MIT License](LICENSE).
