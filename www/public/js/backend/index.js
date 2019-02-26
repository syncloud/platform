backend.apps_data = {
    "apps": [
        {
            "id": "wordpress",
            "name": "WordPress",
            "icon": "appsimages/penguin.png",
            "url": "http://owncloud.odroid-c2.syncloud.it"
        }
    ]
};

//    apps_data: {
//      "apps": []
//    },

backend.installed_apps = function (on_complete, on_error) {
    var that = this;
    setTimeout(function () {
        on_complete(that.apps_data);
    }, 2000);
};
