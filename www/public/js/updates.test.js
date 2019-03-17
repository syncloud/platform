const Settings = require('./updates');
const mock = require('../__mocks__/jquery.mockjax');

test( "settings check version", () => {

  $.ajaxSetup({ async: false });
  var actualData;
  Settings.check_versions(function(data) {
    actualData = data;
  }, function(a, b, c) {});

  expect(actualData).isDefined();

});

test( "settings platform upgrade", () => {

  $.ajaxSetup({ async: false });
  Settings.platform_upgrade(function(data) {}, function(a, b, c) {});

});

test( "settings sam upgrade", () => {

  $.ajaxSetup({ async: false });
  Settings.sam_upgrade(function(data) {}, function(a, b, c) {});

});
