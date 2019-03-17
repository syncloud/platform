import App from './app.js'

test( "app install", () => {

  App.run_app_action('owncloud', 'install', function() {}, function(a, b, c) {});

});
