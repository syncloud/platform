import * as App from './app.js'

test( "app install", () => {

  App.run_app_action('owncloud', 'install', function() {}, function(a, b, c) {});

});

test( "backup job status predicate is running", () => {
  const response = {data: 'JobStatusBusy'}
  const is_running = App.backup_status_predicate(response);
  expect(is_running).toEqual(true);
});

test( "backup job status predicate is not running", () => {
  const response = {data: 'JobStatusIdle'}
  const is_running = App.backup_status_predicate(response);
  expect(is_running).toEqual(false);
});

