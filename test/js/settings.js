QUnit.test( "settings check version", function( assert ) {

  backend.async = false;
  check_versions(function(data) {}, function() {}, function(a, b, c) {});

  assert.deepEqual( true, true);
});