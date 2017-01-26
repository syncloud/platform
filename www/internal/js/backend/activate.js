var backend = {
    activate: function(parameters) {
        var values = parameters.values;
        $.post("/rest/activate", values)
            .done(function (data) {
                if (parameters.hasOwnProperty("done")) {
                    parameters.done(data);
                }
            })
            .fail(function (xhr, textStatus, errorThrown) {
                var error = null;
                if (xhr.hasOwnProperty('responseJSON')) {
                    var error = xhr.responseJSON;
                }
                if (parameters.hasOwnProperty("fail")) {
                    parameters.fail(xhr.status, error);
                }
            })
            .always(function() {
                if (parameters.hasOwnProperty("always")) {
                    parameters.always();
                }
            });
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