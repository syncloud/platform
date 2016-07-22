$(document).ready(function () {
    $("#btn_send_logs").on('click', function () {
        var btn = $(this);
        btn.button('loading');

        $.get('/rest/user')
                .done(function() {
                    $.get('/rest/send_log')
                            .done(function (data) {
                                window.location.href = "/";
                            })
                            .fail(onError)
                            .always(function() {
                                btn.button('reset');
                            });
                })
                .fail(function() {
                    window.location.href = "/sendlogs.html";
                });
    });
});