const { available_apps } = require('./appcenter');
require('./backend.mock/common');
require('./backend.mock/appcenter');

test( "show apps success", () => {

  backend.async = false;
  backend.available_apps_success = true
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  available_apps(function(data) {}, on_error);

  expect(error_counter).toEqual(0);

});

test( "show apps error", () => {

  backend.async = false;
  backend.available_apps_success = false
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  available_apps(function(data) {}, on_error);

  expect(error_counter).toEqual(1);

});
