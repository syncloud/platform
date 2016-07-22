function show_error(error) {
    if ('parameters_messages' in error) {
        for (var i = 0; i < error.parameters_messages.length; i++) {
            var pm = error.parameters_messages[i];
            var message_text = pm.messages.join('\n');
        }
    } else {
        if (!('message' in error && error.message)) {
            error['message'] = 'Server Error'
        }
    }
}

$(document).ready(function () {

    $("#form-logs").submit(function (event) {
        event.preventDefault();

        var btn = $('#btn_send');
        btn.button('loading');

        var values = $("#form-logs").serializeArray();
        $.post("/rest/send_log", values)
                .done(function (data) {
                    window.location.href = "/activate.html";
                })
                .fail(function (xhr, textStatus, errorThrown) {
                    if (xhr.hasOwnProperty('responseJSON')) {
                        var error = xhr.responseJSON;
                        show_error(error);
                    } else {
                        window.location.href = "/sendlogs.html";
                    }
                })
                .always(function() {
                    btn.button('reset');
                });
    });

});
