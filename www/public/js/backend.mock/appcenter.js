backend.apps_data = {
      "apps": [
        {
          "id": "owncloud",
          "name": "ownCloud",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "diaspora",
          "name": "Diaspora",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "mail",
          "name": "Mail",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "talk",
          "name": "Talk",
          "icon": "appsimages/penguin.png"
        },
        {
          "id": "files",
          "name": "Files Browser",
          "icon": "appsimages/penguin.png"
        }
      ]
    };

backend.apps_data_error = {
      "message": "error",
      "success": false
    };

backend.available_apps_success = true;

backend.available_apps = function(on_complete, on_error) {
    var that = this;
    if (backend.available_apps_success) {
        backend.test_timeout(function() { on_complete(that.apps_data); }, 2000);
    } else {
        backend.test_timeout(function() { on_complete(that.apps_data_error); }, 2000);
    }
};