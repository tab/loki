local framework = require("framework")
local auth = {}

local token_cache = {
  admin = nil,
  manager = nil,
  user = nil
}

function auth.get_admin_token()
  if token_cache.admin then
    print("Using cached admin token")
    return token_cache.admin
  end

  local token = framework.authenticate_with_smart_id("EE", "40504040001")
  if not token then
      error("Failed to get admin token")
  end

  token_cache.admin = token
  return token
end

function auth.get_manager_token()
  if token_cache.manager then
    print("Using cached manager token")
    return token_cache.manager
  end

  local token = framework.authenticate_with_smart_id("BE", "00010299944")
  if not token then
      error("Failed to get manager token")
  end

  token_cache.manager = token
  return token
end

function auth.get_user_token()
  if token_cache.user then
    print("Using cached user token")
    return token_cache.user
  end

  local token = framework.authenticate_with_smart_id("EE", "30303039914")
  if not token then
      error("Failed to get user token")
  end

  token_cache.user = token
  return token
end

function auth.get_invalid_token()
  return "invalid-token"
end

return auth
