const Settings = require('./settings');
require('./backend.mock/common');
require('./backend.mock/settings');

test( "settings check version", () => {

  backend.async = false;
  Settings.check_versions(function(data) {}, function(a, b, c) {});

});

test( "settings platform upgrade", () => {

  backend.async = false;
  Settings.platform_upgrade(function(data) {}, function(a, b, c) {});

});

test( "settings sam upgrade", () => {

  backend.async = false;
  Settings.sam_upgrade(function(data) {}, function(a, b, c) {});

});