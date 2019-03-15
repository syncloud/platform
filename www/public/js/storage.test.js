const Storage = require('./storage');
const mock = require('../__mocks__/jquery.mockjax');

test( "settings boot disk extend", () => {

  $.ajaxSetup({ async: false });
  Storage.boot_extend(function(data) {}, function(a, b, c) {});

});

test( "settings disk activate error", () => {

  $.ajaxSetup({ async: false });
  mock.disk_action_success = false
  var error_counter = 0;
  function on_error(a, b, c) {
    error_counter += 1
  }

  Storage.disk_action('device', true, function() {}, on_error)

  expect(error_counter).toEqual(1);
});

test( "settings disk activate success", () => {

  $.ajaxSetup({ async: false });
  mock.disk_action_success = true
  var success_counter = 0;
  function on_success(a) {
    success_counter += 1
  }

  Storage.disk_action('device', true, on_success, function() {});

  expect(success_counter).toEqual(1);
});

