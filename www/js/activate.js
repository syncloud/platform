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

function reset_error() {
    $('#errors_placeholder').empty();
    $('#group-redirect-email').removeClass('has-error');
    $('#help-redirect-email').text('');
    $('#group-redirect-password').removeClass('has-error');
    $('#help-redirect-password').text('');
}

function show_error(error) {
    if ('parameters_messages' in error) {
        for (var i = 0; i < error.parameters_messages.length; i++) {
            var pm = error.parameters_messages[i];
            var group_id = '#group-redirect-' + pm.parameter;
            $(group_id).addClass('has-error');
            var hint_id = '#help-redirect-' + pm.parameter;
            var message_text = pm.messages.join('\n');
            $(hint_id).text(message_text);
        }
    } else {
        if (!('message' in error && error.message)) {
            error['message'] = 'Server Error'
        }
        $('#errors_placeholder').append(_.template($("#error_template").html())(error));
    }
}