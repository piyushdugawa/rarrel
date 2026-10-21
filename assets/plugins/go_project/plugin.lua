print("Go Project Lua Plugin Loaded!")

Name = "go-project"
Version = "0.0.1"

ProjectName = prompt("Enter Projcet Name:")

ModulePath = prompt("Enter Module Path:")

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

