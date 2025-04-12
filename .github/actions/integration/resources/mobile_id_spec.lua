local framework = require("framework")
local json = require("cjson")

local suite = {}

local test_data = {
  success_user = { phone_number = "+37268000769", personal_code = "60001017869" },
  not_mid_client = { phone_number = "+37200000266", personal_code = "60001019939" },
  delivery_error = { phone_number = "+37207110066", personal_code = "60001019947" },
  user_cancelled = { phone_number = "+37201100266", personal_code = "60001019950" },
  signature_hash_mismatch = { phone_number = "+37200000666", personal_code = "60001019961" },
  sim_error = { phone_number = "+37201200266", personal_code = "60001019972" },
  phone_absent = { phone_number = "+37213100266", personal_code = "60001019983" },
  timeout = { phone_number = "+37266000266", personal_code = "50001018908" }
}

function suite.test_success_authentication()
  local user = test_data.success_user
  print(string.format("Testing successful authentication with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local token = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code)
  if not token then
    error("Failed to authenticate with valid credentials")
  end

  framework.assert.not_equals(nil, token, "Token should not be nil")
  framework.assert.not_equals("", token, "Token should not be empty")

  return token
end

function suite.test_not_mid_client()
  local user = test_data.not_mid_client
  print(string.format("Testing NOT_MID_CLIENT scenario with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local result = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code, "NOT_MID_CLIENT")
  framework.assert.equals(true, result, "Authentication with expected NOT_MID_CLIENT error should succeed")

  return true
end

function suite.test_delivery_error()
  local user = test_data.delivery_error
  print(string.format("Testing DELIVERY_ERROR scenario with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local result = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code, "DELIVERY_ERROR")
  framework.assert.equals(true, result, "Authentication with expected DELIVERY_ERROR error should succeed")

  return true
end

function suite.test_user_cancelled()
  local user = test_data.user_cancelled
  print(string.format("Testing USER_CANCELLED scenario with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local result = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code, "USER_CANCELLED")
  framework.assert.equals(true, result, "Authentication with expected USER_CANCELLED error should succeed")

  return true
end

function suite.test_signature_hash_mismatch()
  local user = test_data.signature_hash_mismatch
  print(string.format("Testing SIGNATURE_HASH_MISMATCH scenario with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local result = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code, "SIGNATURE_HASH_MISMATCH")
  framework.assert.equals(true, result, "Authentication with expected SIGNATURE_HASH_MISMATCH error should succeed")

  return true
end

function suite.test_sim_error()
  local user = test_data.sim_error
  print(string.format("Testing SIM_ERROR scenario with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local result = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code, "SIM_ERROR")
  framework.assert.equals(true, result, "Authentication with expected SIM_ERROR error should succeed")

  return true
end

function suite.test_phone_absent()
  local user = test_data.phone_absent
  print(string.format("Testing PHONE_ABSENT scenario with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local result = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code, "PHONE_ABSENT")
  framework.assert.equals(true, result, "Authentication with expected PHONE_ABSENT error should succeed")

  return true
end

function suite.test_timeout()
  local user = test_data.timeout
  print(string.format("Testing TIMEOUT scenario with Mobile-ID (Phone: %s, Personal Code: %s)", user.phone_number, user.personal_code))

  local result = framework.authenticate_with_mobile_id(user.phone_number, user.personal_code, "TIMEOUT")
  framework.assert.equals(true, result, "Authentication with expected TIMEOUT error should succeed")

  return true
end

function suite.run()
  local tests = {
    suite.test_success_authentication,
    suite.test_not_mid_client,
    suite.test_delivery_error,
    suite.test_user_cancelled,
    suite.test_signature_hash_mismatch,
    suite.test_sim_error,
    suite.test_phone_absent,
    suite.test_timeout
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
