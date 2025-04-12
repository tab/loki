local framework = require("framework")
local json = require("cjson")

local suite = {}

local test_data = {
  success_user = { country = "EE", personal_code = "40504040001" },
  user_refused = { country = "EE", personal_code = "30403039917" },
  user_refused_display_text_and_pin = { country = "EE", personal_code = "30403039928" },
  user_refused_vc_choice = { country = "EE", personal_code = "30403039939" },
  user_refused_confirmation_message = { country = "EE", personal_code = "30403039946" },
  user_refused_confirmation_message_with_vc_choice = { country = "EE", personal_code = "30403039950" },
  user_refused_cert_choice = { country = "EE", personal_code = "30403039961" },
  wrong_vc = { country = "EE", personal_code = "30403039972" }
}

function suite.test_success_authentication()
  local user = test_data.success_user
  print(string.format("Testing successful authentication with Smart-ID (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local token = framework.authenticate_with_smart_id(user.country, user.personal_code)
  if not token then
    error("Failed to authenticate with valid credentials")
  end

  framework.assert.not_equals(nil, token, "Token should not be nil")
  framework.assert.not_equals("", token, "Token should not be empty")

  return token
end

function suite.test_user_refused()
  local user = test_data.user_refused
  print(string.format("Testing USER_REFUSED scenario with Smart-ID (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local result = framework.authenticate_with_smart_id(user.country, user.personal_code, "USER_REFUSED")
  framework.assert.equals(true, result, "Authentication with expected USER_REFUSED error should succeed")

  return true
end

function suite.test_user_refused_display_text_and_pin()
  local user = test_data.user_refused_display_text_and_pin
  print(string.format("Testing USER_REFUSED_DISPLAYTEXTANDPIN scenario (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local result = framework.authenticate_with_smart_id(user.country, user.personal_code, "USER_REFUSED_DISPLAYTEXTANDPIN")
  framework.assert.equals(true, result, "Authentication with expected USER_REFUSED_DISPLAYTEXTANDPIN error should succeed")

  return true
end

function suite.test_user_refused_vc_choice()
  local user = test_data.user_refused_vc_choice
  print(string.format("Testing USER_REFUSED_VC_CHOICE scenario (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local result = framework.authenticate_with_smart_id(user.country, user.personal_code, "USER_REFUSED_VC_CHOICE")
  framework.assert.equals(true, result, "Authentication with expected USER_REFUSED_VC_CHOICE error should succeed")

  return true
end

function suite.test_user_refused_confirmation_message()
  local user = test_data.user_refused_confirmation_message
  print(string.format("Testing USER_REFUSED_CONFIRMATIONMESSAGE scenario (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local result = framework.authenticate_with_smart_id(user.country, user.personal_code, "USER_REFUSED_CONFIRMATIONMESSAGE")
  framework.assert.equals(true, result, "Authentication with expected USER_REFUSED_CONFIRMATIONMESSAGE error should succeed")

  return true
end

function suite.test_user_refused_confirmation_message_with_vc_choice()
  local user = test_data.user_refused_confirmation_message_with_vc_choice
  print(string.format("Testing USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE scenario (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local result = framework.authenticate_with_smart_id(user.country, user.personal_code, "USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE")
  framework.assert.equals(true, result, "Authentication with expected USER_REFUSED_CONFIRMATIONMESSAGE_WITH_VC_CHOICE error should succeed")

  return true
end

function suite.test_user_refused_cert_choice()
  local user = test_data.user_refused_cert_choice
  print(string.format("Testing USER_REFUSED_CERT_CHOICE scenario (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local result = framework.authenticate_with_smart_id(user.country, user.personal_code, "USER_REFUSED_CERT_CHOICE")
  framework.assert.equals(true, result, "Authentication with expected USER_REFUSED_CERT_CHOICE error should succeed")

  return true
end

function suite.test_wrong_vc()
  local user = test_data.wrong_vc
  print(string.format("Testing WRONG_VC scenario (Country: %s, Personal Code: %s)", user.country, user.personal_code))

  local result = framework.authenticate_with_smart_id(user.country, user.personal_code, "WRONG_VC")
  framework.assert.equals(true, result, "Authentication with expected WRONG_VC error should succeed")

  return true
end

function suite.run()
  local tests = {
    suite.test_success_authentication,
    suite.test_user_refused,
    suite.test_user_refused_display_text_and_pin,
    suite.test_user_refused_vc_choice,
    suite.test_user_refused_confirmation_message,
    suite.test_user_refused_confirmation_message_with_vc_choice,
    suite.test_user_refused_cert_choice,
    suite.test_wrong_vc
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
