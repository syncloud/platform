function backend_send_log(on_completed) {
    setTimeout(function() {
        on_completed();
    }, 2000);
}