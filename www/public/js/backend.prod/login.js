backend.login = function(values, on_complete, on_error, on_always) {
        $.post("/rest/login", values).done(on_complete).fail(on_error).always(on_always);
    };