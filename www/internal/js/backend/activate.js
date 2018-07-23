var backend = {
    activate: function(parameters, on_always, on_done, on_error) {
        $.post("/rest/activate", parameters)
            .done(on_done)
            .fail(on_error)
            .always(on_always);
    },
    activate_custom_domain: function(parameters, on_always, on_done, on_error) {
        $.post("/rest/activate_custom_domain", parameters)
            .done(on_done)
            .fail(on_error)
            .always(on_always);
    },
    login: function() {
        var url = (new URI())
                .protocol('https')
                .port(443)
                .filename("")
                .query("");

        window.location.href = url;

    }
};