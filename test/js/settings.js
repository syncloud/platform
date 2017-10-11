QUnit.test( "settings check version", function( assert ) {

  backend.async = false;
  check_versions(function(data) {}, function(a, b, c) {});

  assert.deepEqual( true, true);
});

QUnit.test( "settings platform upgrade", function( assert ) {

  backend.async = false;
  platform_upgrade(function(data) {}, function(a, b, c) {});

  assert.deepEqual( true, true);
});

QUnit.test( "settings sam upgrade", function( assert ) {

  backend.async = false;
  sam_upgrade(function(data) {}, function(a, b, c) {});

  assert.deepEqual( true, true);
});

QUnit.test( "settings boot disk extend", function( assert ) {

  backend.async = false;
  boot_extend(function(data) {}, function(a, b, c) {});

  assert.deepEqual( true, true);
});

QUnit.test( "settings disk activate error", function( assert ) {

  backend.async = false;
  backend.disk_action_success = false
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  disk_action('device', true, function() {}, on_error)

  assert.deepEqual( error_counter, 1);
});

QUnit.test( "settings disk activate success", function( assert ) {

  backend.async = false;
  backend.disk_action_success = true
  var success_counter = 0;
  function on_success(a) {
    success_counter += 1
  }

  disk_action('device', true, on_success, function() {});

  assert.deepEqual( success_counter, 1);
});