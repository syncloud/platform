import * as Updates from './updates.js'
import { Setup } from '../__mocks__/jquery.mockjax.js'

test( "updates check version", () => {

  $.ajaxSetup({ async: false });
  var actualData;
  Updates.check_versions(function(data) {
    actualData = data;
  }, function(a, b, c) {});

  expect(actualData).toBeDefined();

});

test( "updates platform upgrade", () => {

  $.ajaxSetup({ async: false });
  var actualData;
  Updates.platform_upgrade(function(data) {
    actualData = data;
  }, function(a, b, c) {});

  expect(actualData).toBeDefined();

});

test( "updates installer upgrade", () => {

  $.ajaxSetup({ async: false });
  var actualData;
  Updates.installer_upgrade(function(data) {
    actualData = data;
  }, function(a, b, c) {});

  expect(actualData).toBeDefined();

});
