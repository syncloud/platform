backend.device_data = {
      "device_url": "http://test.syncloud.it",
      "success": true
    };

backend.access_data = {
      "data": {
        "external_access": true,
        "protocol": "https"
      },
      "success": true
    };

backend.versions_data = {
      "data": [
        {
          "app": {
            "id": "platform",
            "name": "Platform",
            "required": true,
            "ui": false,
            "url": "http://platform.odroid-c2.syncloud.it"
          },
          "current_version": "880",
          "installed_version": "876"
        },
        {
          "app": {
            "id": "sam",
            "name": "Syncloud App Manager",
            "required": true,
            "ui": false,
            "url": "http://sam.odroid-c2.syncloud.it"
          },
          "current_version": "78",
          "installed_version": "75"
        }
      ],
      "success": true
    };

backend.device_url = function(on_complete, on_error) {
        var that = this;
        backend.test_timeout(function() { on_complete(that.device_data); }, 2000);
    };

backend.send_logs = function(include_support, on_always, on_error) {
        backend.test_timeout(on_always, 2000);
    };

backend.reactivate = function(on_complete, on_error) {
        on_complete({ activate_url: "../internal/activate.html", success: true});
    };

backend.get_versions = function(on_complete, on_error) {
        var that = this;
        backend.test_timeout(function() {
            on_complete(that.versions_data);
        }, 2000);
    };

backend.check_versions = function(on_always, on_error) {
        backend.test_timeout(on_always, 2000);
    };

backend.platform_upgrade = function(on_complete, on_error) {
        backend.test_timeout(function() { on_complete({success: true}) }, 2000);
    };

backend.sam_upgrade = function(on_complete, on_error) {
        backend.test_timeout(function() { on_complete({success: true}) }, 2000);
    };
