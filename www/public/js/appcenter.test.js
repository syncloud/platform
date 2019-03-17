const AppCenter = require('./appcenter');
const mock = require('../__mocks__/jquery.mockjax');


test( "show apps success", () => {

  $.ajaxSetup({ async: false });
  mock.available_apps_success = true;
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  AppCenter.available_apps(function(data) {}, on_error);

  expect(error_counter).toEqual(0);

});

test( "show apps error", () => {

  $.ajaxSetup({ async: false });
  mock.available_apps_success = false;
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  AppCenter.available_apps(function(data) {}, on_error);

  expect(error_counter).toEqual(1);

});
