local R = {}

function R.chk_tools(list)
    for index, dictionary in ipairs(list) do
        local tool = dictionary["name"]
        local version = dictionary["version"]
        local cmd = dictionary["check_command"]

        local path, exists = exist(tool)

        if exists then
            print("woohooo! " .. tool .. " exists at: " .. path)
            if version == "" or version ~= nil then
                print("Found: " .. extractVer)
            else
                if extractVer(cmd) == version or extractVer(cmd) > version then
                    print("Version satisfied!!!")
                else
                    print("Version not satisfied!!!")
                end
            end
            
        else
            print("Ahhhh! " .. tool .. " is not installed or not in PATH")
        end
    end
end

return R
