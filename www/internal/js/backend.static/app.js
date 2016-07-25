var backend = {
    app_data: {
      "info": {
        "app": {
          "id": "owncloud",
          "name": "ownCloud",
          "required": false,
          "icon": "penguin.png",
          "ui": true,
          "url": "/"
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