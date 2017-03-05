var backend = {
    send_log: function(on_always, on_error) {
        $.get('/rest/send_log').always(on_always).fail(on_error);
    }
};