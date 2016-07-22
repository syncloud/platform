function auto_login(name, password) {
    var url = (new URI())
            .port(80)
            .directory("rest")
            .filename("login")
            .query("");

    var form = $(
            '<form action="' + url + '" method="post">' +
            '<input type="hidden" name="name" value="' + name + '" />' +
            '<input type="hidden" name="password" value="' + password + '" />' +
            '</form>');
    $('body').append(form);
    form.submit();
}

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
    $("#form_activate").submit(function (event) {
        event.preventDefault();

        var btn = $('#btn_activate');
        btn.button('loading');

        var values = $("#form_activate").serializeArray();
        $.post("/rest/activate", values)
                .done(function (data) {
                    var device_username = $('#device_username').val();
                    var device_password = $('#device_password').val();
                    auto_login(device_username, device_password)
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