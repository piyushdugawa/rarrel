local R = {}

function R.chk_tools(list)
    for index, dictionary in ipairs(list) do
        local tool = dictionary["name"]
        exist(tool)
    end
end

return R
