var app_data = {
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

function backend_app_action(app_id, action, on_complete) {
    setTimeout(function() {
        backend_update_app(app_id, on_complete);
    }, 2000);
}

function backend_update_app(app_id, on_complete) {
    setTimeout(function() {
        ui_display_app(app_data);
        on_complete();
    }, 2000);
}