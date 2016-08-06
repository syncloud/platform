var backend = {
    app_data: {
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
    },
    load_app: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.app_data);
        }, 2000);
s   },
    app_action: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters);
        }, 2000);
    }
}