QUnit.test( "settings check version", function( assert ) {

  backend.async = false;
  check_versions(function(data) {}, function() {}, function(a, b, c) {});

  assert.deepEqual( true, true);
});

QUnit.test( "settings platform upgrade", function( assert ) {

  backend.async = false;
  platform_upgrade(function(data) {}, function() {}, function(a, b, c) {});

  assert.deepEqual( true, true);
});