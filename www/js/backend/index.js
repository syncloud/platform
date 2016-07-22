function backend_installed_apps(on_completed) {
    $.get( '/rest/installed_apps')
        .done( function(data) {
            on_completed(data);
        })
        .fail( onError );
}