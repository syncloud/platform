backend.get_versions = function(on_complete, on_error) {
        $.get('/rest/settings/versions').done(on_complete).fail(on_error);
    };

backend.check_versions = function(on_always, on_error) {
        $.get('/rest/check').always(on_always).fail(on_error);
    };

backend.platform_upgrade = function(on_complete, on_error) {
        $.get('/rest/upgrade', { app_id: 'platform' }).done(on_complete).fail(on_error);
    };

backend.sam_upgrade = function(on_complete, on_error) {
        $.get('/rest/upgrade', { app_id: 'sam' }).done(on_complete).fail(on_error);
    };
