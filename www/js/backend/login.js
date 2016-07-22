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