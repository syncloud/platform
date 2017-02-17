var backend = {
    apps_data: {
      "apps": [
        {
          "id": "owncloud",
          "name": "ownCloud",
          "icon": "penguin.png"
        },
        {
          "id": "diaspora",
          "name": "Diaspora",
          "icon": "penguin.png"
        },
        {
          "id": "mail",
          "name": "Mail",
          "icon": "penguin.png"
        },
        {
          "id": "talk",
          "name": "Talk",
          "icon": "penguin.png"
        },
        {
          "id": "files",
          "name": "Files Browser",
          "icon": "penguin.png"
        }
      ]
    },

    available_apps: function(parameters) {
        var that = this;
        setTimeout(function() {
            success_callbacks(parameters, that.apps_data);
        }, 2000);
    }
}