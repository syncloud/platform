var backend = {
    login: function(parameters) {
        var values = parameters.values;
        $.post("/rest/login", values)
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
    }
}

function backend_login(values, on_error) {
    $.post("/rest/login", values)
            .done(function (data) {
                window.location.replace("/");
            })
            .fail(function (xhr, textStatus, errorThrown) {
                if (xhr.hasOwnProperty('responseJSON')) {
                    var error = xhr.responseJSON;
            		typeof on_error === 'function' && on_error(error);
                } else {
                    window.location.href = "login.html";
                }
            });
}