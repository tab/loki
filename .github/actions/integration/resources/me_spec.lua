local framework = require("framework")
local auth = require("auth")
local json = require("cjson")

local suite = {}

local test_data = {
  admin_token = nil,
  manager_token = nil,
  user_token = nil,
}

function suite.setup()
  print("Setting up tests...")
  test_data.admin_token = auth.get_admin_token()
  test_data.manager_token = auth.get_manager_token()
  test_data.user_token = auth.get_user_token()

  return test_data.admin_token ~= nil
end

function suite.test_get_admin()
  print("Test: Get me with admin token")

  local response = framework.request(
    "GET",
    framework.config.loki_url .. "/api/me",
    {
      ["Authorization"] = "Bearer " .. test_data.admin_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )
  framework.assert.status_code(response, 200)

  print("Test: Get me with user token")

  local response = framework.request(
    "GET",
    framework.config.loki_url .. "/api/me",
    {
      ["Authorization"] = "Bearer " .. test_data.user_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 200)

  return true
end

function suite.test_get_manager()
  print("Test: Get me with manager token")

  local response = framework.request(
    "GET",
    framework.config.loki_url .. "/api/me",
    {
      ["Authorization"] = "Bearer " .. test_data.manager_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 200)

  return true
end

function suite.test_get_user()
  print("Test: Get me with user token")

  local response = framework.request(
    "GET",
    framework.config.loki_url .. "/api/me",
    {
      ["Authorization"] = "Bearer " .. test_data.user_token,
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 200)

  return true
end

function suite.test_unauthorized()
  print("Test: Unauthorized access")

  local response = framework.request(
    "GET",
    framework.config.loki_url .. "/api/me",
    {
      ["Authorization"] = "Bearer " .. auth.get_invalid_token(),
      ["X-Trace-ID"] = framework.generate_trace_id(),
      ["X-Request-ID"] = framework.generate_request_id()
    }
  )

  framework.assert.status_code(response, 401)

  return true
end

function suite.run()
  if not suite.setup() then
    print("❌ Setup failed")
    return false
  end

  local tests = {
    suite.test_get_admin,
    suite.test_get_manager,
    suite.test_get_user,
    suite.test_unauthorized
  }

  local success = true

  for i, test in ipairs(tests) do
    local test_success, result = pcall(test)

    if not test_success or not result then
      print("❌ Test failed: " .. debug.traceback())
      success = false
      break
    else
      print("✅ Test passed")
    end
  end

  return success
end

return suite
