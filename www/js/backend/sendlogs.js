function backend_send_log(values, on_done, on_error, on_always) {
    $.post("/rest/send_log", values)
        .done(function (data) {
            on_done();
        })
        .fail(function (xhr, textStatus, errorThrown) {
            var error = null;
            if (xhr.hasOwnProperty('responseJSON')) {
                error = xhr.responseJSON;
           }
            on_error(error);
        })
        .always(function() {
            on_always();
        });
}