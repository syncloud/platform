var apps_data = {
  "apps": [
    {
      "id": "owncloud",
      "name": "ownCloud",
      "icon": "penguin.png",
      "url": "http://owncloud.odroid-c2.syncloud.it"
    }
  ]
};

function backend_installed_apps(on_completed) {
    setTimeout(function() {
        on_completed(apps_data);
    }, 2000);
}