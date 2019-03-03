const backwnd = require('./settings');
require('./backend.mock/common');
require('./backend.mock/settings');

test( "settings check version", () => {

  backend.async = false;
  check_versions(function(data) {}, function(a, b, c) {});

});

test( "settings platform upgrade", () => {

  backend.async = false;
  platform_upgrade(function(data) {}, function(a, b, c) {});

});

test( "settings sam upgrade", () => {

  backend.async = false;
  sam_upgrade(function(data) {}, function(a, b, c) {});

});

test( "settings boot disk extend", () => {

  backend.async = false;
  boot_extend(function(data) {}, function(a, b, c) {});

});

test( "settings disk activate error", () => {

  backend.async = false;
  backend.disk_action_success = false
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  disk_action('device', true, function() {}, on_error)

  expect(error_counter).toEqual(1);
});

test( "settings disk activate success", () => {

  backend.async = false;
  backend.disk_action_success = true
  var success_counter = 0;
  function on_success(a) {
    success_counter += 1
  }

  disk_action('device', true, on_success, function() {});

  expect(success_counter).toEqual(1);
});

