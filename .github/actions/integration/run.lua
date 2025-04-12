local framework = require("framework")
local mobile_id_suite = require("resources/mobile_id_spec")
local smart_id_suite = require("resources/smart_id_spec")
local me_suite = require("resources/me_spec")

if arg[1] == "--debug" then
  framework.debug(true)
  print("Debug mode enabled")
end

local results = {
  start_time = os.time(),
  tests = {},
  passed = 0,
  failed = 0,
  suite_test_count = 0,
  suite_tests_passed = 0,
  suite_tests_failed = 0
}

local function format_duration(seconds)
  if seconds < 60 then
    return string.format("%.2f seconds", seconds)
  else
    local minutes = math.floor(seconds / 60)
    local remaining_seconds = seconds % 60
    return string.format("%d minutes and %.2f seconds", minutes, remaining_seconds)
  end
end

local function count_tests_in_suite(suite)
  local count = 0
  for name, func in pairs(suite) do
    if type(func) == "function" and name:match("^test_") then
      count = count + 1
    end
  end
  return count
end

local function run_test_suite(name, suite)
  print("\n" .. string.rep("=", 80))
  print("Running " .. name .. " tests...")
  print(string.rep("=", 80))

  local suite_test_count = count_tests_in_suite(suite)
  results.suite_test_count = results.suite_test_count + suite_test_count

  local test_start = os.time()
  local success = suite.run()
  local test_end = os.time()
  local duration = test_end - test_start

  table.insert(results.tests, {
    name = name,
    success = success,
    result = success and "PASSED" or "FAILED",
    duration = duration,
    test_count = suite_test_count
  })

  if success then
    results.passed = results.passed + 1
    results.suite_tests_passed = results.suite_tests_passed + suite_test_count
    print("\n✅ Test suite passed in " .. format_duration(duration) .. " (" .. suite_test_count .. " tests)")
  else
    results.failed = results.failed + 1
    results.suite_tests_failed = results.suite_tests_failed + suite_test_count
    print("\n❌ Test suite failed after " .. format_duration(duration) .. " (" .. suite_test_count .. " tests)")
  end
end

local function check_services()
  print("Checking if services are up...")

  local response = framework.loki_readiness()
  if not response or response.status ~= 200 then
    print("❌ Loki service is not ready")
    return false
  end

  local response = framework.loki_backoffice_readiness()
  if not response or response.status ~= 200 then
    print("❌ Loki-backoffice service is not ready")
    return false
  end

  print("✅ All services are ready")
  return true
end

if check_services() then
  run_test_suite("Mobile-ID", mobile_id_suite)
  run_test_suite("Smart-ID", smart_id_suite)
  run_test_suite("Me", me_suite)
else
  print("❌ Services are not ready. Skipping tests.")
  os.exit(1)
end

results.end_time = os.time()
results.total_duration = results.end_time - results.start_time

print("\n" .. string.rep("=", 80))
print("INTEGRATION TEST RESULTS")
print(string.rep("=", 80))
print("Total duration: " .. format_duration(results.total_duration))
print("Test suites: " .. results.passed .. " passed, " .. results.failed .. " failed, " .. #results.tests .. " total")
print("Tests:       " .. results.suite_tests_passed .. " passed, " .. results.suite_tests_failed .. " failed, " .. results.suite_test_count .. " total")
print(string.rep("=", 80))

for _, test in ipairs(results.tests) do
  local status_icon = test.success and "✅" or "❌"
  print(status_icon .. " " .. test.name .. ": " .. test.result .. " (" .. test.test_count .. " tests in " .. format_duration(test.duration) .. ")")
end

os.exit(results.failed > 0 and 1 or 0)
