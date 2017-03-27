backend.app_data = {
      "info": {
        "app": {
          "id": "owncloud",
          "name": "ownCloud",
          "required": false,
          "icon": "penguin.png",
          "ui": true,
          "url": "/",
          "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborumLorem ipsum dolor sit amet, consectetur"
        },
        "current_version": "212",
        "installed_version": "210"
      }
    };
    
backend.load_app = function(app_id, on_complete, on_error) {
        var that = this;
        backend.test_timeout(function() { on_complete(that.app_data); }, 2000);
    };
   
backend.app_action = function(app_id, action, on_complete, on_error) {
        var that = this;
        backend.test_timeout(function() {on_complete({success: true, is_running: false});}, 2000);
    };
