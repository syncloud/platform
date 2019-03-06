backend.device_url = function(on_complete, on_error) {
        $.get('/rest/settings/device_url').done(on_complete).fail(on_error);
    };

backend.send_logs = function(include_support, on_always, on_error) {
        $.get('/rest/send_log',
          { include_support: include_support }
        ).always(on_always).fail(on_error);
    };

backend.reactivate = function(on_complete, on_error) {
        $.get('/rest/settings/activate_url').done(on_complete).fail(on_error);
    };