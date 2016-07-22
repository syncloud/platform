function backend_send_log(on_completed) {
    $.get('/rest/user')
        .done(function() {
            $.get('/rest/send_log')
                .done(function (data) {
                    window.location.href = "/";
                })
                .fail(onError)
                .always(function() {
                    on_completed();
                });
        })
        .fail(function() {
            window.location.href = "/sendlogs.html";
        });
}