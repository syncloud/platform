function backend_send_log(values, on_done, on_error, on_always) {
    setTimeout(function() {
        on_done();
        on_always();
    }, 2000);
}