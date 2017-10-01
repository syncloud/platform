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
    login: function(name, password) {
        var url = (new URI())
                .port(80)
                .filename("/rest/login")
                .query("");

        var form = $(
                '<form action="' + url + '" method="post">' +
                '<input type="hidden" name="name" value="' + name + '" />' +
                '<input type="hidden" name="password" value="' + password + '" />' +
                '</form>');
        $('body').append(form);
        form.submit();
    }
};