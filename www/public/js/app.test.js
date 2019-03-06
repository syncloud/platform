const App = require('./app');
require('./backend.mock/common');
require('./backend.mock/app');

test( "app install", () => {

  backend.async = false;
  App.run_app_action('owncloud', 'install', function() {}, function(a, b, c) {});

});
