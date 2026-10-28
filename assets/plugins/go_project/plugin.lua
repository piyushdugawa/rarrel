print("\nGo Project Lua Plugin Loaded!\n")

Name = "go-project"
Version = "0.0.1"

local pwd = io.popen("pwd"):read("*l")
local basename = pwd:match("([^/]+)/*$")

ProjectName = prompt("Enter Project Name:", "? ", basename)
print("\27[34m◆\27[0m \27[1;97mProject Name:\27[0m \27[38;5;201m" .. ProjectName .. "\27[0m")

local dest

if ProjectName == basename then
    dest = "./"
else
    dest = ProjectName
end

ModulePath = prompt("Enter Module Path:", "? ", "github.com/piyushdugawa/rarrel")
print("\27[34m◆\27[0m \27[1;97mModule:\27[0m \27[38;5;201m" .. ModulePath .. "\27[0m")


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
        name = "jq",
        required = false,
        version_arg = "--version"
    }
}

GoVersion = Tools[1].version

verifyTools(Tools)

renderTemplates("./templates", dest)

