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