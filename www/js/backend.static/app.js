$( document ).ready(function () {
    var data = {
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
    };
    display_app(data);
});
