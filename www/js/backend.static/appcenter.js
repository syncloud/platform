$( document ).ready(function() {
    var data = {
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
    };
    display_apps(data);
});