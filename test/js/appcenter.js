QUnit.test( "show apps success", function( assert ) {


  backend.available_apps_success = true
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  available_apps(function(data) {}, on_error);

  assert.deepEqual( error_counter, 0);

});

QUnit.test( "show apps error", function( assert ) {


  backend.available_apps_success = false
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  available_apps(function(data) {}, on_error);

  assert.deepEqual( error_counter, 1);

});
