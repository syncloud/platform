function send_logs() {
    var btn = $(this);
    btn.prop('disabled', true);
    btn.text('Sending ...');
    $.get('/rest/send_log')
            .done(function () {
                btn.text('Send logs');
                btn.prop('disabled', false);
            }).fail(onError);
}