function auto_login(name, password) {
    var url = (new URI())
            .port(80)
            .directory("server/rest")
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

function reset_error() {
    $('#errors_placeholder').empty();
    $('#group-email').removeClass('has-error');
    $('#help-email').text('');
    $('#group-password').removeClass('has-error');
    $('#help-password').text('');
}

function show_error(error) {
    if ('parameters_messages' in error) {
        for (var i = 0; i < error.parameters_messages.length; i++) {
            var pm = error.parameters_messages[i];
            var group_id = '#group-' + pm.parameter;
            $(group_id).addClass('has-error');
            var hint_id = '#help-' + pm.parameter;
            var message_text = pm.messages.join('\n');
            $(hint_id).text(message_text);
        }
    } else {
        $('#errors_placeholder').append(_.template($("#error_template").html())(error));
    }
}