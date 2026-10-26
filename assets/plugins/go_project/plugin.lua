print("\nGo Project Lua Plugin Loaded!\n")

Name = "go-project"
Version = "0.0.1"

ProjectName = prompt("Enter Project Name:", "? ", "rrl")
print("\27[34m◆\27[0m \27[1;97mProject Name:\27[0m \27[38;5;201m" .. ProjectName .. "\27[0m")

ModulePath = prompt("Enter Module Path:", "? ", "rarrel")
print("\27[34m◆\27[0m \27[1;97mModule:\27[0m \27[38;5;201m" .. ModulePath .. "\27[0m")

Tools = {
    {
        name = "go",
        required = true,
        version = ">=1.24"
    },
    {
        name = "git",
        required = false,
        version = "inherit"
    }
}

Templates = "./templates"

