var backend_menu = {
    logout: function (on_complete, on_error) {
        $.post('/rest/logout').done(on_complete).fail(on_error);
    },
    restart: function (on_complete, on_error) {
        $.get('/rest/restart').done(on_complete).fail(on_error);
    },
    shutdown: function (on_complete, on_error) {
        $.get('/rest/shutdown').done(on_complete).fail(on_error);
    }
};