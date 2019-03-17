const Settings = require('./updates');
const mock = require('../__mocks__/jquery.mockjax');

test( "settings check version", () => {

  $.ajaxSetup({ async: false });
  var actualData;
  Settings.check_versions(function(data) {
    actualData = data;
  }, function(a, b, c) {});

  expect(actualData).toBeDefined();

});

test( "settings platform upgrade", () => {

  $.ajaxSetup({ async: false });
  var actualData;
  Settings.platform_upgrade(function(data) {
    actualData = data;
  }, function(a, b, c) {});

  expect(actualData).toBeDefined();

});

test( "settings sam upgrade", () => {

  $.ajaxSetup({ async: false });
  var actualData;
  Settings.sam_upgrade(function(data) {
    actualData = data;
  }, function(a, b, c) {});

  expect(actualData).toBeDefined();

});
